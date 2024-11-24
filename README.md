# RepoRover

RepoRover is a CLI tool for managing multiple repositories through group-focused operations, providing commands to initialize groups, add or remove repositories, perform bulk updates, and customize configurations.

## Features

- **Group Management**: Organize repositories into groups for collective operations.
- **Bulk Operations**: Execute commands like `status`, `pull`, and `sync` across all repositories in a group.
- **Customization**: Set global and group-specific configurations.
- **Templates**: Use templates to quickly set up new groups with predefined settings.
- **Aliases**: Create shortcuts for frequently used commands.

## Installation

To install RepoRover, you can use `pip`:

```bash
<install command(s)>
```

## Quick Start

### 1. Initialize a New Group

Create a new group to organize your repositories:

```bash
rover group init <group name>
```

### 2. Add Repositories to the Group

Add repositories by providing their URLs or local paths:

```bash
rover group add <group name> https://github.com/user/repo1.git ~/projects/repo2
```

### 3. View Group Details

List all repositories within a group:

```bash
rover group show <group name>
```

### 4. Pull Updates for All Repositories

Update all repositories in the group:

```bash
rover group pull <group name>
```

### 5. Check the Status of Repositories

Get the status of all repositories in the group:

```bash
rover group status <group name>
```

### 6. Execute a Command Across the Group

Run a custom command on all repositories:

```bash
rover group exec <group name> -- "git fetch --all"
```

## Examples

- **Sync Repositories to a Branch:**

  ```bash
  rover group sync <group name> --branch develop
  ```

- **Analyze Commits Since a Date:**

  ```bash
  rover group analyze commits <group name> --since="2023-01-01"
  ```

- **Remove a Repository from a Group:**

  ```bash
  rover group remove <group name> repo-name
  ```

- **Set a Group-Specific Configuration:**

  ```bash
  rover group config set <group name> default-branch develop
  ```

## License

RepoRover is released under the [MIT License](LICENSE).