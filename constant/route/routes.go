package route

const BasePath = "/api/v1"

const UserPath = BasePath + "/users"
const UserLogin = UserPath + "/login"
const UserLoginGoogle = UserPath + "/login-google"
const UserGoogleCallback = UserPath + "/google-callback"
const UserRegister = UserPath + "/register"
const UserUpdateAvatar = UserPath + "/avatar"

const AdminPath = BasePath + "/admin"
const AdminLogin = AdminPath + "/login"
const AdminEdit = AdminPath + "/edit/:id"
const AdminDelete = AdminPath + "/delete"

const AdminManageUserPath = AdminPath + "/users"
const AdminManageUserByID = AdminManageUserPath + "/:id"

const ProductPath = BasePath + "/products"
const CategoryProduct = ProductPath + "/categories/:category_name"
const ProductByID = ProductPath + "/:id"

const ImpactCategoryPath = BasePath + "/impacts"
const ImpactCategoryByID = ImpactCategoryPath + "/:id"
