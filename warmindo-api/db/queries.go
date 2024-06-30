package db

const (
	CreateCategoryQuery  = `INSERT INTO categories (name) VALUES ($1)`
	GetCategoriesQuery   = `SELECT id, name, created_at, updated_at FROM categories`
	GetCategoryByIDQuery = `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`
	UpdateCategoryQuery  = `UPDATE categories SET name = $1, updated_at = NOW() WHERE id = $2`
	DeleteCategoryQuery  = `DELETE FROM categories WHERE id = $1`

	CreateMenuQuery  = `INSERT INTO menus (name, image, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`
	GetMenusQuery    = `SELECT id, name, image, description, price, category_id, created_at, updated_at FROM menus`
	GetMenuByIDQuery = `SELECT id, name, image, description, price, category_id, created_at, updated_at FROM menus WHERE id = $1`
	UpdateMenuQuery  = `UPDATE menus SET name = $1, image = $2, description = $3, price = $4, category_id = $5, updated_at = NOW() WHERE id = $6`
	DeleteMenuQuery  = `DELETE FROM menus WHERE id = $1`

	CreateOrderQuery  = `INSERT INTO orders (amount, table_number, status_id, order_date, menu_id) VALUES ($1, $2, $3, $4, $5)`
	GetOrdersQuery    = `SELECT id, amount, table_number, status_id, order_date, menu_id, created_at, updated_at FROM orders`
	GetOrderByIDQuery = `SELECT id, amount, table_number, status_id, order_date, menu_id, created_at, updated_at FROM orders WHERE id = $1`
	UpdateOrderQuery  = `UPDATE orders SET amount = $1, table_number = $2, status_id = $3, order_date = $4, menu_id = $5, updated_at = NOW() WHERE id = $6`
	DeleteOrderQuery  = `DELETE FROM orders WHERE id = $1`

	CreateRoleQuery  = `INSERT INTO roles (name) VALUES ($1)`
	GetRolesQuery    = `SELECT id, name, created_at, updated_at FROM roles`
	GetRoleByIDQuery = `SELECT id, name, created_at, updated_at FROM roles WHERE id = $1`
	UpdateRoleQuery  = `UPDATE roles SET name = $1, updated_at = NOW() WHERE id = $2`
	DeleteRoleQuery  = `DELETE FROM roles WHERE id = $1`

	CreateStatusQuery  = `INSERT INTO statuses (name) VALUES ($1)`
	GetStatusesQuery   = `SELECT id, name, created_at, updated_at FROM statuses`
	GetStatusByIDQuery = `SELECT id, name, created_at, updated_at FROM statuses WHERE id = $1`
	UpdateStatusQuery  = `UPDATE statuses SET name = $1, updated_at = NOW() WHERE id = $2`
	DeleteStatusQuery  = `DELETE FROM statuses WHERE id = $1`

	CreateUserQuery  = `INSERT INTO users (email, password, name, username, role_id, phone) VALUES ($1, $2, $3, $4, $5, $6)`
	GetUsersQuery    = `SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM users`
	GetUserByIDQuery = `SELECT id, email, name, username, role_id, phone, created_at, updated_at FROM users WHERE id = $1`
	UpdateUserQuery  = `UPDATE users SET email = $1, password = $2, name = $3, username = $4, role_id = $5, phone = $6, updated_at = NOW() WHERE id = $7`
	DeleteUserQuery  = `DELETE FROM users WHERE id = $1`
)
