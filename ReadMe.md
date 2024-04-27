# Simple CLI Password Manager

The Simple CLI Password Manager is a command-line tool that allows users to securely store and manage their passwords. It provides a convenient way to generate strong passwords, store them in an encrypted format, and retrieve them when needed.

## Features

- **Password Storage**: Safely store passwords for various accounts and services.
- **Password Generation**: Generate strong, randomized passwords with customizable length and character sets.
- **Encryption**: Encrypt password data to ensure security and privacy.
- **Search**: Easily search for stored passwords by account name or keyword.
- **Update and Delete**: Modify or remove existing passwords as needed.
- **Master Password**: Protect the password database with a master password for added security.

## Installation

To install the Simple CLI Password Manager, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/your_username/simple-cli-password-manager.git
    ```

2. Navigate to the project directory:

    ```bash
    cd simple-cli-password-manager
    ```

3. Build the executable:

    ```bash
    go build -o password-manager
    ```

4. Run the executable:

    ```bash
    ./password-manager
    ```

## Usage

After running the executable, you'll be prompted to create a master password for the password manager. Once set up, you can use the following commands:

- **add**: Add a new password entry.
- **get**: Retrieve a password for a specific account.
- **list**: List all stored passwords.
- **update**: Update an existing password entry.
- **delete**: Delete a password entry.
- **generate**: Generate a new random password.
- **search**: Search for passwords containing a specific keyword.
- **exit**: Exit the password manager.

Here's an example of how to add a new password entry:

```bash
./password-manager add
```

Follow the prompts to enter the account name, username, password, and any additional notes.

## Security Considerations

- **Master Password**: Choose a strong and unique master password to protect your password database.
- **Encryption**: Password data is encrypted using industry-standard encryption algorithms to prevent unauthorized access.
- **Avoid Storing Sensitive Information**: While the password manager encrypts data, it's best to avoid storing highly sensitive information such as financial or personal data.

## Contributing

Contributions are welcome! If you encounter any bugs or have suggestions for improvements, please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the [MIT License](LICENSE).
