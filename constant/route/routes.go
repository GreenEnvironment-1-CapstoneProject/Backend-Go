package route

const BasePath = "/api/v1"

const UserPath = BasePath + "/user"
const UserLogin = UserPath + "/login"
const UserRegister = UserPath + "/register"

const AdminPath = BasePath + "/admin"
const AdminLogin = AdminPath + "/login"
const AdminEdit = AdminPath + "/edit/:id"
const AdminDelete = AdminPath + "/delete"
