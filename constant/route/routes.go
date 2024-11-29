package route

const BasePath = "/api/v1"

const UserPath = BasePath + "/users"
const UserLogin = UserPath + "/login"
const UserLoginGoogle = UserPath + "/login-google"
const UserGoogleCallback = UserPath + "/google-callback"
const UserRegister = UserPath + "/register"

const AdminPath = BasePath + "/admin"
const AdminLogin = AdminPath + "/login"
const AdminEdit = AdminPath + "/edit/:id"
const AdminDelete = AdminPath + "/delete"
