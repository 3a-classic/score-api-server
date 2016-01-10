#curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostScoreEntrySheet localhost:8443/api/page/scoreViewSheet/A/
#curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostApply localhost:8443/api/page/applyScore/A/


curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostLogin.json test.3a-classic.com/v1/page/login

## userRegister
#curl -X POST -H "Content-Type : application/json"  -d @testUserRegister.json test.3a-classic.com/v1/register/user

## teamRegister
#curl -X POST -H "Content-Type : application/json"  -d @testTeamRegister.json test.3a-classic.com/v1/register/20151121/team

## fieldRegister
#curl -X POST -H "Content-Type : application/json"  -d @testFieldRegister.json test.3a-classic.com/v1/register/20151121/field
