<p align="center">
    <h1 align="center">CULINARY-BLISS</h1>
</p>
<p align="center">
    <em>Culinary Bliss, a estaurant management app designed to streamline your operations and elevate your dining experience. Whether you own a bustling city bistro or a cozy countryside café, Culinary Bliss is your go-to solution for efficient, effective, and effortless restaurant management.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/ShahSau/culinary-bliss?style=flat&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/ShahSau/culinary-bliss?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/ShahSau/culinary-bliss?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/ShahSau/culinary-bliss?style=flat&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
		<em>Developed with the software and tools below.</em>
</p>
<p align="center">
		<img src="https://img.shields.io/badge/YAML-CB171E.svg?style=flat&logo=YAML&logoColor=white" alt="YAML">
	<img src="https://img.shields.io/badge/JSON-000000.svg?style=flat&logo=JSON&logoColor=white" alt="JSON">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
	<img src="https://img.shields.io/badge/JWT-000000?style=flat&logo=Go&logoColor=white" alt="JWT">
        <img src="https://img.shields.io/badge/Gin-black?style=flat&logo=Go&logoColor=white" alt="Gin-go">
</p>
<hr>

##  Quick Links

> - [📍 Overview](#-overview)
> - [📦 Features](#-features)
> - [📂 Repository Structure](#-repository-structure)
> - [🧩 Modules](#-modules)
> - [🚀 Getting Started](#-getting-started)
>   - [⚙️ Installation](#️-installation)
>   - [🤖 Running EthnicElegance](#-running-EthnicElegance)
>   - [🧪 Tests](#-tests)
> - [🛠 Project Roadmap](#-project-roadmap)
> - [📄 License](#-license)
> - [👏 Acknowledgments](#-acknowledgments)

---

##  Overview

HTTP error 401 for prompt `overview`

---

##  Features

This restaurant management backend, built with Go and Gin, provides robust and scalable functionalities to support a full-fledged restaurant management. The backend includes the following key features:

<h6>JWT Authentication</h6>
<p>Secure Authentication: Implemented JWT (JSON Web Token) for secure and stateless user authentication.</p>
<p>Token Generation: Generate tokens for user sessions upon successful login.</p>
<p>Token Verification: Verify tokens for protected routes to ensure only authenticated users can access certain endpoints.</p>
<p>Refresh Tokens: Support for refreshing tokens to maintain secure sessions.</p>
  
<h6>Admin Section</h6>
<p>User Management: Create, read, update, and delete (CRUD) operations for managing user accounts.</p>
<p>Role Management: Assign and manage roles (e.g., admin, staff) to control access levels and permissions.</p>
<p>Restaurant Management: Manage restaurant details, including name, location, hours of operation, and more.</p>
  
<h6>Global Section</h6>
<p>Menu Management: CRUD operations for managing the restaurant menu, including categories, items, prices, and availability.</p>
<p>Reservation Management: Handle customer reservations, including booking, updating, and canceling reservations.</p>
<p>Table Management: Manage table assignments, statuses, and seating arrangements to optimize dining space.</p>
  
<h6>User Section</h6>
<p>User Registration and Login: Allow users to register and log in securely.</p>
<p>Profile Management: Users can update their profiles, including personal details and preferences.</p>
<p>Reservation Booking: Users can book, view, and cancel reservations.</p>

---

##  Repository Structure

```sh
└── culinary-bliss/
    ├── LICENSE
    ├── README.md
    ├── controllers
    │   ├── authControllers.go
    │   ├── categeoryControllers.go
    │   ├── foodControllers.go
    │   ├── invoiceControllers.go
    │   ├── menuControllers.go
    │   ├── orderControllers.go
    │   ├── orderItemControllers.go
    │   ├── restaurantControllers.go
    │   ├── tableControllers.go
    │   └── userControllers.go
    ├── culinary-bliss
    ├── database
    │   └── databaseConnection.go
    ├── docs
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    ├── go.mod
    ├── go.sum
    ├── helpers
    │   ├── adminHelper.go
    │   └── tokenHelper.go
    ├── main.go
    ├── middleware
    │   └── authMiddleware.go
    ├── models
    │   ├── categoryModel.go
    │   ├── foodModel.go
    │   ├── invoiceModel.go
    │   ├── menuModel.go
    │   ├── orderItemModel.go
    │   ├── orderModel.go
    │   ├── restaurantModel.go
    │   ├── tableModel.go
    │   └── userModel.go
    ├── routes
    │   ├── authRouter.go
    │   ├── catgeoryRouter.go
    │   ├── foodRouter.go
    │   ├── globalRouter.go
    │   ├── invoiceRouter.go
    │   ├── menuRouter.go
    │   ├── orderItemRouter.go
    │   ├── orderRouter.go
    │   ├── restaurantRouter.go
    │   ├── tableRouter.go
    │   └── userRouter.go
    └── types
        ├── category-type.go
        ├── invoice-type.go
        ├── menu-type.go
        ├── restaurant-type.go
        ├── table-type.go
        └── user-type.go
```

---

##  Modules


<details closed><summary>helpers</summary>

| File                                                                                           | Summary                                            |
| ---                                                                                            | ---                                                |
| [adminHelper.go](https://github.com/ShahSau/culinary-bliss/blob/master/helpers/adminHelper.go) | Admin Helper functions |
| [tokenHelper.go](https://github.com/ShahSau/culinary-bliss/blob/master/helpers/tokenHelper.go) | Token helper functions |

</details>

<details closed><summary>middleware</summary>

| File                                                                                                    | Summary                                                  |
| ---                                                                                                     | ---                                                      |
| [authMiddleware.go](https://github.com/ShahSau/culinary-bliss/blob/master/middleware/authMiddleware.go) | authentication middleware |

</details>

<details closed><summary>types</summary>

| File                                                                                                 | Summary                                              |
| ---                                                                                                  | ---                                                  |
| [category-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/category-type.go)     | category type  |
| [menu-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/menu-type.go)             | menu type       |
| [table-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/table-type.go)           | table type.      |
| [invoice-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/invoice-type.go)       | invoice type    |
| [restaurant-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/restaurant-type.go) | restaurant type |
| [user-type.go](https://github.com/ShahSau/culinary-bliss/blob/master/types/user-type.go)             | user type      |

</details>

<details closed><summary>database</summary>

| File                                                                                                          | Summary                                                    |
| ---                                                                                                           | ---                                                        |
| [databaseConnection.go](https://github.com/ShahSau/culinary-bliss/blob/master/database/databaseConnection.go) | databaseConnection |

</details>

<details closed><summary>models</summary>

| File                                                                                                  | Summary                                               |
| ---                                                                                                   | ---                                                   |
| [invoiceModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/invoiceModel.go)       | invoiceModel    |
| [foodModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/foodModel.go)             | foodModel       |
| [menuModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/menuModel.go)             | menuModel      |
| [orderModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/orderModel.go)           | orderModel      |
| [orderItemModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/orderItemModel.go)   | orderItemModel  |
| [tableModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/tableModel.go)           | tableModel.     |
| [restaurantModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/restaurantModel.go) | restaurantModel |
| [categoryModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/categoryModel.go)     | categoryModel   |
| [userModel.go](https://github.com/ShahSau/culinary-bliss/blob/master/models/userModel.go)             | userModel       |

</details>

<details closed><summary>routes</summary>

| File                                                                                                    | Summary                                                |
| ---                                                                                                     | ---                                                    |
| [foodRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/foodRouter.go)             | foodRouter      |
| [menuRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/menuRouter.go)             | menuRouter      |
| [authRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/authRouter.go)             | authRouter       |
| [tableRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/tableRouter.go)           | tableRouter      |
| [globalRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/globalRouter.go)         |globalRouter     |
| [invoiceRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/invoiceRouter.go)       |invoiceRouter    |
| [orderItemRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/orderItemRouter.go)   |orderItemRouter  |
| [orderRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/orderRouter.go)           |orderRouter      |
| [userRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/userRouter.go)             |userRouter      |
| [catgeoryRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/catgeoryRouter.go)     |catgeoryRouter  |
| [restaurantRouter.go](https://github.com/ShahSau/culinary-bliss/blob/master/routes/restaurantRouter.go) |restaurantRouter|

</details>

<details closed><summary>controllers</summary>

| File                                                                                                                   | Summary                                                          |
| ---                                                                                                                    | ---                                                              |
| [orderControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/orderControllers.go)           | orderControllers      |
| [menuControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/menuControllers.go)             | menuControllers       |
| [restaurantControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/restaurantControllers.go) | restaurantControllers |
| [categeoryControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/categeoryControllers.go)   | categeoryControllers  |
| [orderItemControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/orderItemControllers.go)   | orderItemControllers  |
| [tableControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/tableControllers.go)           | tableControllers      |
| [authControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/authControllers.go)             | authControllers       |
| [foodControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/foodControllers.go)             | foodControllers       |
| [userControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/userControllers.go)             | userControllers      |
| [invoiceControllers.go](https://github.com/ShahSau/culinary-bliss/blob/master/controllers/invoiceControllers.go)       | invoiceControllers    |

</details>

---

##  Getting Started

***Requirements***

Ensure you have the following dependencies installed on your system:

* **Go**: `version 1.22.5`

###  Installation

1. Clone the culinary-bliss repository:

```sh
git clone https://github.com/ShahSau/culinary-bliss
```

2. Change to the project directory:

```sh
cd culinary-bliss
```

3. Install the dependencies:

```sh
go build -o myapp
```

###  Running culinary-bliss

Use the following command to run culinary-bliss:

```sh
./myapp
```

###  Tests

To execute tests, run:

```sh
go test
```

---

##  Project Roadmap

- [ ] `► API Testing`
- [ ] `► Docker`
- [ ] `► CI/CD`

---

##  License

This project is protected under the MIT License. For more details, refer to the [LICENSE](https://github.com/ShahSau/turbo?tab=MIT-1-ov-file#readme) file.

---

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.

[**Return**](#-quick-links)

---
