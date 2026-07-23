-- Let the dev user create and drop the throwaway database used by `go test`
-- (see pkg/db/testing.go). Runs as root on first container startup.
GRANT ALL PRIVILEGES ON `zauth_test`.* TO 'zauth'@'%';
FLUSH PRIVILEGES;
