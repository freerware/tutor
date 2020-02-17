/*
	Grant appropriate privileges for web_app user.
*/
-- REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'app_user'@'%'; -- since web_app is created using $MYSQL_USER and $MYSQL_PASSWORD, it has GRANT ALL applied which isn't good.
-- GRANT SELECT, INSERT, UPDATE, DELETE on safest-watch.* TO 'app_user'@'%'; -- adding back the privileges that are more restrictive.
-- FLUSH PRIVILEGES;
