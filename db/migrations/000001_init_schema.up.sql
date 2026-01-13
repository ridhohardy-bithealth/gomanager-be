-- Create table users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    user_image_uri VARCHAR(255),
    company_name VARCHAR(255),
    company_image_uri VARCHAR(255)
);

-- Create table departments
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create enum
CREATE TYPE enum_gender as ENUM ('male', 'female');

-- Create employees
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender enum_gender,
    identity_number VARCHAR(255) NOT NULL,
    employee_image_uri VARCHAR(255),
    department_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT employees_identity_number_per_user UNIQUE (user_id, identity_number)
);