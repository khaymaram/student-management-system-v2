# Student Management System

A full-stack university administration application for managing students, professors, courses, enrollments, schedules, majors, grades, and finances. The application includes role-based authentication for administrators, students, and professors.

## Technology stack

- Frontend: React 19, TypeScript, Vite, Tailwind CSS, Radix UI, TanStack Query, Axios, and Zod
- Backend: Go, Gin, GORM, and bcrypt
- Database: Microsoft SQL Server
- Authentication: university ID and password with a signed JWT session

## Project structure

```text
student-management-system-v2/
├── sms-frontend/        React web application
│   └── src/
│       ├── components/  Pages, layouts, and reusable UI components
│       ├── context/     Authentication state
│       ├── hooks/       API queries and mutations
│       ├── lib/         Axios client and utilities
│       └── types/       TypeScript and Zod schemas
└── sms-backend/         Go REST API
    ├── cmd/             Server, migration, and seed entry points
    └── internal/
        ├── handlers/    HTTP request handling
        ├── middleware/  Authentication and authorization
        ├── models/      Database entities
        ├── repositories/ Database access
        ├── services/    Business rules
        └── routes/      API route definitions
```

## Prerequisites

- Node.js and npm
- Go 1.26.4 or a compatible version
- Microsoft SQL Server
- A local SQL Server database named `test`
- Windows integrated authentication access to that database

The current database connector uses this local development connection directly:

```text
server=localhost;database=test;integrated security=true
```

Although database environment variables are represented in the configuration package, the current connector does not use them. Update `sms-backend/internal/initializers/database.go` if a different SQL Server host, database, or authentication method is required.

## Initial setup

### 1. Configure the backend

Create `sms-backend/.env` if it does not already exist:

```dotenv
SERVER_PORT=8081
JWT_SECRET=replace-with-a-long-random-secret
```

Always set a stable `JWT_SECRET`. If it is empty, the backend generates a new secret at startup and all existing sessions become invalid whenever the server restarts.

### 2. Install frontend dependencies

From the repository root:

```powershell
cd sms-frontend
npm install
```

The frontend defaults to `http://localhost:8081/api`. To use another API address, create `sms-frontend/.env`:

```dotenv
VITE_API_URL=http://localhost:8081/api
```

### 3. Create and seed the database

From `sms-backend`:

```powershell
go run ./cmd/migrate
go run ./cmd/seed
```

> Warning: the migration command currently drops and recreates all application tables. It deletes existing users, majors, students, finances, professors, courses, and enrollments. Do not run it against data that must be preserved.

### 4. Start the application

Start the backend from `sms-backend`:

```powershell
go run ./cmd/server
```

Start the frontend in a separate terminal from `sms-frontend`:

```powershell
npm run dev
```

Open `http://localhost:5173`.

On Windows systems where PowerShell blocks `npm.ps1`, use `npm.cmd` instead:

```powershell
npm.cmd run dev
```

## Seeded logins

These accounts are created by `go run ./cmd/seed`:

| Role | University ID | Password |
|---|---|---|
| Administrator | `ADMIN001` | `Admin123!` |
| Student | `1122` | `Student1122!` |
| Professor | `P1001` | `TeacherP1001!` |

The seed passwords are intended for local development only and should be changed before using the application with real data.

## Automatically created accounts

When an administrator creates a student or professor, a linked login account is created automatically.

| Account | Default password | Generated email pattern |
|---|---|---|
| Student | `Student[student ID]!` | `[normalized name][student ID]@grgi.edu` |
| Professor | `Teacher[professor ID]!` | `[normalized name][professor ID]@grgi.edu` |

For example, student `1200` receives password `Student1200!`. Users sign in with their student or professor ID, not their email address. Passwords are stored as bcrypt hashes.

## Roles and permissions

| Feature | Administrator | Student | Professor |
|---|:---:|:---:|:---:|
| System dashboard | Yes | No | No |
| Manage students and professors | Yes | No | No |
| Manage courses, majors, and finances | Yes | No | No |
| View own courses and schedule | N/A | Yes | Yes |
| Enroll in available courses | Yes | Yes | No |
| View and pay own fees | N/A | Yes | No |
| View students in taught courses | Yes | No | Yes |
| Grade students in taught courses | Yes | No | Yes |
| Edit own name, email, and password | Yes | Yes | Yes |
| Delete own account | Admin only | No | No |

Deleting a student or professor from the administrator interface also deletes the associated login account. Updating a name from the User Account page updates both the login profile and the linked student or professor record.

## Core business rules

### Majors

- Students must select a major from the database-backed major list.
- `Undeclared` is included as a valid option.
- Administrators can add and delete majors from the Majors tab.
- `Undeclared` cannot be deleted.
- Deleting another major automatically reassigns its students to `Undeclared` in the same database transaction.

### Courses and scheduling

- Meeting days are Monday through Friday.
- Start times are available on the hour or half-hour.
- The first class starts at 8:00 AM and the last selectable start is 4:00 PM.
- Every class lasts one hour.
- A professor cannot be assigned overlapping courses.
- A student cannot enroll in overlapping courses.
- Rescheduling a course is rejected if it conflicts with its professor or an enrolled student.
- A student cannot enroll in more than 15 credits per semester.

### Grades

- Course grades use `A`, `B`, `C`, `D`, or `F`.
- Student GPA is recalculated from graded enrollments and course credits.
- Professors can grade students only in courses they teach.

### Finances

- In-state tuition is `$12,000`.
- Out-of-state tuition is `$25,000`.
- Residency is selected when creating a student and can be changed from the Edit Student or Finance interface.
- Students may pay only their own outstanding balance.
- Students cannot change their scholarship or residency from the payment request.
- Payments cannot reduce the amount previously paid or exceed the remaining balance.

## API overview

The backend serves JSON endpoints under `/api`. With the default configuration, the base URL is `http://localhost:8081/api`.

Main route groups:

| Route group | Purpose |
|---|---|
| `/auth` | Login, current user, profile, and password management |
| `/students` | Student roster, filters, details, and enrollments |
| `/professors` | Professor records and details |
| `/courses` | Course records, schedules, and rosters |
| `/enrollments` | Enrollment data |
| `/finances` | Tuition, scholarships, residency, and payments |
| `/majors` | Available majors |

Except for `POST /api/auth/login`, API routes require an `Authorization: Bearer <token>` header. Authorization middleware limits each role to its permitted resources.

The health check is available without authentication:

```http
GET /ping
```

## Development commands

Frontend commands, run from `sms-frontend`:

```powershell
npm run dev       # Start the development server
npm run build     # Type-check and build for production
npm run lint      # Run ESLint
npm run preview   # Preview the production build
```

Backend commands, run from `sms-backend`:

```powershell
go run ./cmd/server   # Start the API
go run ./cmd/migrate  # Drop and recreate the schema
go run ./cmd/seed     # Insert local sample data
go test ./...         # Run backend tests
```

## Troubleshooting

### Login says invalid

- Confirm the university ID and password match exactly; IDs are used instead of email addresses.
- If seed credentials were changed in code after the database was seeded, the stored bcrypt hash does not change automatically. Recreate and reseed the local database, or use the password originally stored for that account.
- Confirm the backend is running and the frontend points to the correct `VITE_API_URL`.

### Student or professor is not found

- Confirm the login account has a `subject_id` matching an existing Student or Professor record.
- Create linked accounts through the administrator UI so both records are created together.

### Database connection fails

- Confirm SQL Server is running on `localhost`.
- Confirm the `test` database exists.
- Confirm the Windows account running the backend has database access through integrated security.

### Browser requests are blocked by CORS

The backend currently allows the development frontend origin `http://localhost:5173`. Update the CORS configuration in `sms-backend/internal/routes/api.go` when hosting the frontend elsewhere.

## Production considerations

Before deploying this project:

- Replace the destructive development migration with versioned, non-destructive migrations.
- Move the SQL Server connection string fully into environment configuration.
- Use a strong, stable JWT secret managed outside source control.
- Replace all seeded passwords and development accounts.
- Configure production CORS origins.
- Serve both applications over HTTPS.
- Add rate limiting and login-attempt protection.
- Review payment handling before connecting it to any real payment processor.
