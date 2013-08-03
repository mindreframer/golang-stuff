Roles

1. Developer and App

2. Object

Clients:

1. Dashboard: A token and encryption token is shared between dashboard and hamster server to allow developer object creation. On creating developer, developer_id and access_token is released for further developer level actions. access_token is a time based token. The access_token can be stored in a cookie by the dashboard client. A new access_token can be retrieved by logging in. App is on the same access level as developer needing application_id and access_token. Object level authentication is done by app_token and app_secret. They can be retrieved by querying the relevant app.

2. Others(Android, iOS, Java, Python, RoR): Other clients have Only object level permission. The app_token and app_secret is retrieved from the dashboard and manually embedded in the clients.