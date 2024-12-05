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

const CartPath = BasePath + "/cart"
const CartByID = CartPath + "/:id"

const TransactionPath = BasePath + "/transactions"
const TransactionByID = TransactionPath + "/:id"

const ReviewProduct = BasePath + "/reviews"
const ReviewProductByID = ReviewProduct + "/products/:id"

const ChatbotPath = BasePath + "/chatbots"
const ChatbotPathByID = ChatbotPath + "/:chatID"

const ForumPath = BasePath + "/forums"
const ForumByID = ForumPath + "/:id"
const GetForumByUserID = ForumPath + "/user"

const ForumMessagePath = BasePath + "/message" + "/:id"
const ForumMessage = ForumPath + "/message"
const ForumMessageByID = ForumMessage + "/:id"
