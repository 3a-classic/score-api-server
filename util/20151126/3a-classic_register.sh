## userRegister
#curl -X POST -H "Content-Type : application/json"  -d @3a-classic_user.json api.3a-classic.com/v1/register/user
#
### teamRegister
#curl -X POST -H "Content-Type : application/json"  -d @3a-classic_team.json api.3a-classic.com/v1/register/20151126/team
#
### fieldRegister
#curl -X POST -H "Content-Type : application/json"  -d @3a-classic_field.json api.3a-classic.com/v1/register/20151126/field

# userRegister
curl -X POST -H "Content-Type : application/json"  -d @3a-classic_user.json localhost:8080/v1/register/user

## teamRegister
curl -X POST -H "Content-Type : application/json"  -d @3a-classic_team.json localhost:8080/v1/register/team/20151126

## fieldRegister
curl -X POST -H "Content-Type : application/json"  -d @3a-classic_field.json localhost:8080/v1/register/field/20151126
