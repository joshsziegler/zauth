package ldapserver

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/ansel1/merry"
	"github.com/jmoiron/sqlx"
	"github.com/nmcclain/ldap"
	logging "github.com/op/go-logging"

	mGroup "github.com/joshsziegler/zauth/models/group"
	"github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/models/user2group"
)

// LdapConfig is used to pass required configuration options for the LDAP server
type LdapConfig struct {
	BaseDN   string
	UserOU   string
	GroupOU  string
	ListenTo string
}

// We use a global config, because it should be read-only after initial loading
var config LdapConfig

var log *logging.Logger

// DB is the shared (via connection pooling) database connection; goroutine-safe
var DB *sqlx.DB

// Listen performs setup and runs the LDAP server (blocking)
func Listen(logger *logging.Logger, database *sqlx.DB, config LdapConfig) {
	log = logger
	DB = database
	// Create our LDAP-server
	s := ldap.NewServer()
	// Ask the LDAP server to enforce search filter, attribute limits, size/time
	// limits, search scope, and base DN matching to our handler's returned data
	s.EnforceLDAP = true
	// Register LDAP Bind, Search, and Close function handlers
	handler := mysqlBackend{}
	s.BindFunc("", handler)
	s.SearchFunc("", handler)
	s.CloseFunc("", handler)

	// Start the LDAP server
	log.Infof("LDAP server listening on: %s", config.ListenTo)
	err := s.ListenAndServe(config.ListenTo)
	if err != nil {
		log.Fatal("LDAP Server Failed: ", err.Error())
	}
}

// returns the username from the provided UID
//
// Assumes Format:  uid=Username,UserOU,BaseDN
func getUsernameFromUID(uid string) (username string, err error) {
	parts := strings.Split(uid, ",")
	if len(parts) < 1 {
		return "", merry.New("error finding username in UID")
	}
	username = parts[0]
	if len(username) < 5 || username[0:4] != "uid=" {
		return "", merry.New("error finding username in UID")
	}
	return username[4:], nil
}

// Backend interface for LDAP using MySQL as it's datastore
type mysqlBackend struct {
}

// Bind handles client connections, which may be anonymous, or it will have a
// username and password which will be tested against our database of users.
//
// This will NEVER return a user's password hash.
func (h mysqlBackend) Bind(bindDN, bindPassword string, conn net.Conn) (
	ldap.LDAPResultCode, error) {

	if bindDN == "" && bindPassword == "" {
		// Always allow anonymous binds
		log.Info("LDAP: anonymous bind")
		return ldap.LDAPResultSuccess, nil
	}

	// User is trying to bind as a particular user, so check their password
	username, err := getUsernameFromUID(bindDN)
	if err != nil {
		log.Errorf("LDAP: bind failure: could not parse username from %s (%s)",
			bindDN, err)
		return ldap.LDAPResultInvalidCredentials, nil
	}
	tx, err := DB.Beginx()
	if err != nil {
		log.Errorf("LDAP: error starting transaction during Bind: %s", err)
		return ldap.LDAPResultOperationsError, nil
	}
	err = user.Login(tx, username, bindPassword)
	if err != nil {
		log.Errorf("LDAP: bind failure as %s: %s", username, err)
		err = tx.Commit()
		if err != nil {
			log.Errorf("LDAP: transaction error during Bind: %s", err)
			return ldap.LDAPResultOperationsError, nil
		}
		return ldap.LDAPResultInvalidCredentials, nil
	}
	log.Infof("LDAP: bind success as %s", username)
	err = tx.Commit()
	if err != nil {
		log.Errorf("LDAP: transaction error during Bind: %s", err)
		return ldap.LDAPResultOperationsError, nil
	}
	return ldap.LDAPResultSuccess, nil
}

// Search handles a bound client's search request, with LDAP handling the filter
//
// TODO: Only respond to our base DN? - JZ
func (h mysqlBackend) Search(boundDN string, searchReq ldap.SearchRequest,
	conn net.Conn) (ldap.ServerSearchResult, error) {

	// Get username, assuming there will be no error since they already bound
	username, _ := getUsernameFromUID(boundDN)

	scope := ldap.ScopeMap[searchReq.Scope]
	msg := fmt.Sprintf(
		`LDAP: Search by: "%s" BaseDN: "%s" Scope: "%s" Filter: "%s" Attributes: %+v`,
		username, searchReq.BaseDN, scope, searchReq.Filter, searchReq.Attributes)
	entries, err := h.getAllUsersAndGroups()
	if err != nil {
		log.Errorf(`%s FAILED: "%s"`, msg, err)
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, nil
	}
	log.Info(msg)
	return ldap.ServerSearchResult{entries, []string{}, []ldap.Control{},
		ldap.LDAPResultSuccess}, nil
}

// Close handles client disconnections
func (h mysqlBackend) Close(boundDN string, conn net.Conn) error {
	log.Debug("LDAP: closing connection")
	conn.Close()
	return nil
}

// Returns all users and groups in the database as LDAP entries.
//
// This is meant to be passed to the LDAP library for filtering as needed.
func (h mysqlBackend) getAllUsersAndGroups() (entries []*ldap.Entry, err error) {
	tx, err := DB.Beginx()
	if err != nil {
		err = merry.Append(err, "error starting transaction")
		return
	}
	users, groups, err := user2group.GetAll(tx)
	if err != nil {
		err = tx.Commit()
		if err != nil {
			err = merry.Wrap(err)
		}
		return
	}

	for _, user := range users {
		entries = append(entries, userToLDAPEntry(user))
	}
	for _, group := range groups {
		entries = append(entries, groupToLDAPEntry(group))
	}
	err = tx.Commit()
	if err != nil {
		err = merry.Wrap(err)
	}
	return
}

func userToLDAPEntry(u *user.User) *ldap.Entry {
	return &ldap.Entry{"uid=" + u.Username + "," + config.UserOU +
		config.BaseDN,
		[]*ldap.EntryAttribute{
			{"uid", []string{u.Username}},
			{"cn", []string{u.CommonName()}},
			{"sn", []string{u.LastName}},
			{"givenName", []string{u.FirstName}},
			{"uidNumber", []string{strconv.FormatInt(u.UnixUserID(), 10)}},
			{"gidNumber", []string{strconv.FormatInt(u.UnixGroupID(), 10)}},
			{"mail", []string{u.Email}},
			{"homeDirectory", []string{u.HomeDirectory()}},
			{"objectClass", []string{"top"}},
			{"objectClass", []string{"posixAccount"}},
			{"objectClass", []string{"inetOrgPerson"}},
			{"memberOf", u.Groups},
		}}
}

func groupToLDAPEntry(g *mGroup.Group) *ldap.Entry {
	return &ldap.Entry{"cn=" + g.Name + "," + config.GroupOU +
		config.BaseDN,
		[]*ldap.EntryAttribute{
			{"cn", []string{g.Name}},
			{"gidNumber", []string{strconv.FormatInt(g.UnixGroupID, 10)}},
			{"description", []string{g.Description}},
			{"objectClass", []string{"top"}},
			{"objectClass", []string{"posixGroup"}},
			{"objectClass", []string{"groupOfNames"}},
			{"member", g.Members},
		}}
}
