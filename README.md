# Venue Reservation Web App

This is a Go-based web application that allows venue owners to list their venues and customers to browse and reserve them. The application features user authentication, role-based access control, CSRF protection, secure sessions, and a simple admin dashboard for venue management.

## Features

- User signup and login
- Role-based access: venue owners and customers
- Venue creation, editing, and deletion (owner-only)
- Venue listings and reservations (customer-only)
- Review system
- Reservation management
- CSRF and session protection using middleware

## User Roles

- `Owner` (role ID: 1): Can add, edit, and delete venues.
- `Customer` (role ID: 2): Can browse venues and make/cancel reservations.

## Routes

### Public Routes

| Method | Path                | Description                     |
|--------|---------------------|---------------------------------|
| GET    | `/`                 | Home page                       |
| GET    | `/user/signup`      | Show signup form                |
| POST   | `/user/signup`      | Submit new user registration    |
| GET    | `/user/login`       | Show login form                 |
| POST   | `/user/login`       | Log in user                     |
| POST   | `/user/logout`      | Log out user                    |

### Shared Authenticated Routes

| Method | Path                 | Description                              |
|--------|----------------------|------------------------------------------|
| GET    | `/venue/listing`     | View all venues (any authenticated user) |
| GET    | `/venue/{id}`        | View venue details                       |
| POST   | `/venue/{id}/review` | Submit a review                          |

### Owner-Only Routes (role ID: 1)

| Method | Path                   | Description             |
|--------|------------------------|-------------------------|
| GET    | `/venue/form`          | Show new venue form     |
| POST   | `/venue/add`           | Submit new venue        |
| GET    | `/venue/{id}/edit`     | Edit existing venue     |
| POST   | `/venue/{id}/edit`     | Submit venue update     |
| POST   | `/venue/{id}/delete`   | Delete venue            |

### Customer-Only Routes (role ID: 2)

| Method | Path                               | Description                     |
|--------|------------------------------------|---------------------------------|
| POST   | `/reservation/{id}/create`         | Make a reservation              |
| GET    | `/reservations`                    | View all reservations           |
| GET    | `/reservations/cancelled`          | View cancelled reservations     |
| GET    | `/reservations/update/{id}`        | Show update form for reservation|
| POST   | `/reservations/update/{id}`        | Submit reservation update       |
| POST   | `/reservations/cancel/{id}`        | Cancel reservation              |

## Middleware

The app uses `alice` for chaining middleware. Here’s how they’re organized:

- **Standard Middleware**: 
  - `recoverPanic`: Recovers from panics
  - `logRequest`: Logs each request
  - `secureHeaders`: Adds security headers

- **Dynamic Middleware**:
  - `session.Enable`: Enables session management
  - `loggingMiddleware`: Logs user details
  - `authenticate`: Loads and verifies the user
  - `noSurf`: CSRF protection

- **Protected Routes**: Extend dynamic middleware with `requireAuthentication`

- **Role-Based Middleware**:
  - `requireRole(1)`: Only owners
  - `requireRole(2)`: Only customers