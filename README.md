# GoInitializr

This project is a **Golang-based project initializer** designed to help you quickly scaffold a new Go project with common backend features.

---

## 🚀 How to Run

You can run this project in **two ways**:

### 1. 🐳 Run Locally Using Docker

#### 📋 Requirements

- [Docker](https://www.docker.com/)
- [Postman](https://www.postman.com/) (optional, but recommended)

#### 🧰 Steps

1. **Run Docker**

   Start the project using Docker with the appropriate syntax:

   ```bash
   docker build -t go-initializr .
   docker run -d -e PORT=1323 -p 1323:1323 go-initializr
   ```

2. **Hit the API**

   You can test the API in either of the following ways:

   - **Using Postman**

     - Import the provided Postman collection found in this repository.
     - Run the request named **"Generate TemplateProject"**.
     - Update the `{{url}}` variable to `http://localhost:1323/` (default).
     - Click the **Send** button, then use **Download** to retrieve the generated project.

   - **Using Any API Tool with Download Support**

     Send a `POST` request to:

     ```
     http://localhost:1323/v1/initialize
     ```

     #### 📦 Payload Example

     ```json
     {
       "project_name": "my-project",
       "jwt": true,
       "swagger": true,
       "redis": true,
       "validator": true,
       "db": "postgres", // options: postgres, mysql
       "framework": "echo" // currently only 'echo' is supported
     }
     ```

---

### 2. 🌐 Using Deployment Link *(Coming Soon)*

This option will allow you to generate projects from a hosted web service.

_Stay tuned!_

---