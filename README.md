# 🚗 Car Leasing API

A simple RESTful API built with **Go** and **PostgreSQL**, demonstrating authentication, authorization, and permission-based access control.

---

## 🧠 Overview

This project simulates a **car leasing system**.  
It includes management of **cars**, **users**, **reservations**, and **permissions**, using middleware-based permission checks.

---

## ⚙️ Technologies

- **Go (Golang)** – main API
- **PostgreSQL** – database
- **Gorilla Mux** – HTTP router
- **JWT-based Auth** – authentication and permission handling

---

## 🔐 Permissions Model

Permissions are static and represent allowed actions in the system:

| Permission | Description |
|-------------|-------------|
| `write:cars` | Create, update, or delete cars |
| `read:users` | Read user details |
| `write:users` | Grant or revoke user permissions |
| `write:reservations` | Manage reservations |

- Permissions are **not creatable** through the API.
- They can only be **assigned** or **revoked** by users who already have `write:users`.

---

## 🚀 Endpoints

| Method | Endpoint | Description | Auth | Permission |
|--------|-----------|--------------|-------|-------------|
| **Cars** |  |  |  |  |
| `GET` | `/cars` | Get all cars | ❌ | — |
| `GET` | `/cars/{id}` | Get car by ID | ❌ | — |
| `POST` | `/cars` | Add a car | ✅ | `write:cars` |
| `PUT` | `/cars/{id}` | Replace a car | ✅ | `write:cars` |
| `PATCH` | `/cars/{id}` | Update a car | ✅ | `write:cars` |
| `DELETE` | `/cars/{id}` | Delete a car | ✅ | `write:cars` |
| **Users** |  |  |  |  |
| `POST` | `/users` | Sign up a new user | ❌ | — |
| `GET` | `/users/me` | Get current user info | ✅ | — |
| `GET` | `/users/me/permissions` | Get current user's permissions | ✅ | — |
| `GET` | `/users` | List all users | ✅ | `read:users` |
| `GET` | `/users/{id}` | Get user info | ✅ | `read:users` |
| `GET` | `/users/{id}/permissions` | List user’s permissions | ✅ | `read:users` |
| `PUT` | `/users/{id}/permissions/{perm_id}` | Grant permission | ✅ | `write:users` |
| `DELETE` | `/users/{id}/permissions/{perm_id}` | Revoke permission | ✅ | `write:users` |
| **Reservations** |  |  |  |  |
| `GET` | `/reservations` | List reservations (of logged-in user) | ✅ | — |
| `GET` | `/reservations/{id}` | Get reservation by ID | ✅ | (if not owner) `read:reservations` |
| `POST` | `/reservations` | Create reservation | ✅ | `write:reservations` |
| `PUT` | `/reservations/{id}` | Replace reservation | ✅ | `write:reservations` |
| `PATCH` | `/reservations/{id}` | Update reservation | ✅ | `write:reservations` |
| `DELETE` | `/reservations/{id}` | Delete reservation | ✅ | `write:reservations` |
| **Permissions** |  |  |  |  |
| `GET` | `/permissions` | List all permissions | ❌ | — |
| `GET` | `/permissions/{id}` | Get permission details | ❌ | — |
| **Auth** |  |  |  |  |
| `POST` | `/auth/login` | Log in and receive JWT token | ❌ | — |

---

## 🧩 Auth Flow

1. **User signs up** via `POST /users`.
2. **User logs in** via `POST /auth/login` → receives JWT token.
3. All protected routes require `Authorization: Bearer <token>`
4. Middleware (`authMiddleware.WithPerms`) checks permissions before allowing access.

---
