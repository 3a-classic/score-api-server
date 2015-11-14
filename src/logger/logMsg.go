package logger

const (
	Fatal = "fatal"
	Error = "error"
	Info  = "info"
	Debug = "debug"

	ErrMsg   = "Error Message"
	TraceMsg = "Stack Trace"

	E_Nil         = "This variable is nil"
	E_WrongData   = "This do not correct data"
	E_TooManyData = "There are too many data"
	E_MakeHash    = "Can not make hash string"

	E_M_FindEntireCol   = "Can not find entire colection : mongo"
	E_M_FindCol         = "Can not find colection : mongo"
	E_M_Upsert          = "Can not upsert data : mongo"
	E_M_Insert          = "Can not insert data : mongo"
	E_M_Update          = "Can not update data : mongo"
	E_M_RegisterThread  = "Can not register thread score : mongo"
	E_M_RegisterUser    = "Can not register user"
	E_M_SearchPhotoTask = "Can not search picture task : mongo"

	I_M_GetPage     = "Get page data : mongo"
	I_M_PostPage    = "Post page data : mongo"
	I_M_RegisterCol = "Register collection data : mongo"

	E_R_PostPage    = "Can not post page data : route"
	E_R_RegisterCol = "Can not register collection data : route"
	E_R_Upsert      = "Can not upsert data : route"
	E_R_WriteJSON   = "Can not write JSON : route"
	E_R_PingMsg     = "Can not ping message : route"
	E_R_Upgrader    = "Can not upgrader webdocket : route"
)
