#curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostScoreEntrySheet localhost:8443/api/page/scoreViewSheet/A/
#curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostApply localhost:8443/api/page/applyScore/A/

curl -X POST -H "Content-Type : application/json"  -d @testDataOfPostLogin.json test.3a-classic.com/api/page/login
