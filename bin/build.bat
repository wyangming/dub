go build -o ./db.exe ../src/dub/app/db_app.go

go build -o ./serviceUse.exe ../src/dub/app/service/use_app.go

::web service
::useCenter
xcopy "../src/dub/app/web/manlobby/view" "./webrec/manlobby/view" /y /d /s /e /r /f
go build -o ./manLobby.exe ../src/dub/app/web/man_lobby_app.go

::gate server
go build -o ./gateWeb.exe ../src/dub/app/gate/web_app.go

go build -o ./reg.exe ../src/dub/app/reg_app.go