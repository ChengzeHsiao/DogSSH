# Password Authentication Dependencies

DogSSH supports automatic password authentication for SSH connections. To use this feature, you need to install one of the following tools:

## Option 1: sshpass (Recommended)

### macOS
```bash
# Using Homebrew
brew install hudochenkov/sshpass/sshpass
```

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install sshpass
```

### CentOS/RHEL/Fedora
```bash
# CentOS/RHEL
sudo yum install sshpass

# Fedora
sudo dnf install sshpass
```

## Option 2: expect (Fallback)

### macOS
```bash
# Usually pre-installed, or via Homebrew
brew install expect
```

### Ubuntu/Debian
```bash
sudo apt-get update
sudo apt-get install expect
```

### CentOS/RHEL/Fedora
```bash
# CentOS/RHEL
sudo yum install expect

# Fedora
sudo dnf install expect
```

## How It Works

1. When you add a server with a password, DogSSH encrypts and stores the password locally
2. During SSH connection, DogSSH will:
   - First try using `sshpass` if available
   - Fall back to `expect` if `sshpass` is not found
   - Fall back to normal SSH (key-based or interactive) if neither tool is available

## Security Notes

- Passwords are encrypted using AES-256-GCM before being stored locally
- Password files are created with restrictive permissions (0600)
- No passwords are transmitted over the network in plain text
- DogSSH uses the system's native SSH client for all connections