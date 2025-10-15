# ZRide Git Flow Strategy

## Overview
This document outlines the Git branching strategy for the ZRide project, following a modified Git Flow approach optimized for feature development and team collaboration.

## Branch Structure

### 1. Main Branches

#### `master` (Production Branch)
- **Purpose**: Always reflects production-ready state
- **Protection**: Protected branch, no direct commits
- **Deployment**: Automatically deployed to production
- **Merging**: Only from `release/*` or `hotfix/*` branches via PR

#### `develop` (Integration Branch)  
- **Purpose**: Integration branch for features under development
- **Protection**: Protected branch, no direct commits
- **Merging**: From `feature/*` branches via PR after code review
- **Testing**: All features are tested together here

### 2. Supporting Branches

#### `feature/*` (Feature Branches)
- **Naming**: `feature/task-description` or `feature/TASK-001-auth-service`
- **Purpose**: Develop new features or enhancements
- **Lifetime**: Created from `develop`, merged back to `develop`
- **Example**: `feature/auth-service-implementation`

#### `release/*` (Release Branches)
- **Naming**: `release/v1.0.0` or `release/v1.0.0-rc1`  
- **Purpose**: Prepare new production releases
- **Lifetime**: Created from `develop`, merged to both `master` and `develop`
- **Activities**: Bug fixes, documentation, release preparation

#### `hotfix/*` (Hotfix Branches)
- **Naming**: `hotfix/critical-bug-fix` or `hotfix/v1.0.1`
- **Purpose**: Critical fixes for production issues
- **Lifetime**: Created from `master`, merged to both `master` and `develop`
- **Priority**: Highest priority, immediate deployment

#### `bugfix/*` (Bug Fix Branches)
- **Naming**: `bugfix/issue-description`
- **Purpose**: Non-critical bug fixes during development
- **Lifetime**: Created from `develop`, merged back to `develop`

## Workflow for ZRide Tasks

### Current Tasks and Their Branches

```
master
├── develop
│   ├── feature/complete-auth-service
│   ├── feature/api-gateway-setup  
│   ├── feature/user-service-implementation
│   ├── feature/trip-service-implementation
│   ├── feature/matching-service-ai
│   ├── feature/payment-service-zalopay
│   ├── feature/zalo-miniapp-frontend
│   └── feature/deployment-containerization
```

### Branch Naming Convention

#### For Current Todo Items:
1. `feature/complete-auth-service` - Complete Auth Service implementation
2. `feature/zalo-api-research` - Research Zalo Mini App SDK and APIs  
3. `feature/user-service` - Build user management service
4. `feature/trip-service` - Develop trip management service
5. `feature/api-gateway` - Create API Gateway
6. `feature/ai-matching-service` - Build basic AI matching service
7. `feature/zalo-miniapp-frontend` - Develop Zalo Mini App frontend
8. `feature/payment-service` - Integrate payment service
9. `feature/deployment-setup` - Set up containerization and deployment
10. `feature/testing-monitoring` - Implement testing and monitoring

## Git Flow Commands

### 1. Starting a New Feature

```bash
# Ensure you're on develop and up to date
git checkout develop
git pull origin develop

# Create and checkout new feature branch
git checkout -b feature/task-name

# Start development...
# Make commits with descriptive messages
git add .
git commit -m "feat: implement user authentication logic"

# Push feature branch to remote
git push -u origin feature/task-name
```

### 2. Working on a Feature

```bash
# Regular development cycle
git add .
git commit -m "feat: add JWT token validation"
git push origin feature/task-name

# Keep feature branch updated with develop
git checkout develop
git pull origin develop
git checkout feature/task-name
git merge develop
```

### 3. Completing a Feature

```bash
# Final commit and push
git add .
git commit -m "feat: complete auth service implementation"
git push origin feature/task-name

# Create Pull Request to develop branch
# After PR approval and merge, clean up
git checkout develop
git pull origin develop
git branch -d feature/task-name
git push origin --delete feature/task-name
```

### 4. Release Process

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.0.0

# Finalize release (version bumps, documentation)
git add .
git commit -m "chore: prepare v1.0.0 release"
git push -u origin release/v1.0.0

# Create PR to master for release
# After merge, tag the release
git checkout master
git pull origin master
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# Merge release back to develop
git checkout develop
git merge release/v1.0.0
git push origin develop
```

### 5. Hotfix Process

```bash
# Create hotfix from master
git checkout master
git pull origin master
git checkout -b hotfix/critical-security-fix

# Make the fix
git add .
git commit -m "fix: resolve critical security vulnerability"
git push -u origin hotfix/critical-security-fix

# Create PR to master
# After merge, create new tag
git checkout master  
git pull origin master
git tag -a v1.0.1 -m "Hotfix version 1.0.1"
git push origin v1.0.1

# Merge back to develop
git checkout develop
git merge hotfix/critical-security-fix
git push origin develop
```

## Commit Message Convention

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation only changes
- **style**: Code style changes (formatting, semicolons, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvements
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to build process or auxiliary tools
- **ci**: Changes to CI configuration files and scripts

### Examples
```bash
feat(auth): implement Zalo OAuth integration

Add support for Zalo OAuth login flow including:
- Token validation with Zalo API
- User profile retrieval
- Session management

Closes #123

fix(trip): resolve booking seat calculation

Fixed issue where available seats weren't properly updated
when bookings were cancelled.

Fixes #456

docs: add API documentation for auth endpoints

- Added OpenAPI specs for login/logout endpoints
- Updated README with authentication flow
- Added code examples for frontend integration

test(user): add unit tests for user repository

- Added tests for CRUD operations
- Added tests for business rule validation
- Achieved 95% code coverage for user domain
```

## Pull Request Template

### PR Title Format
```
[Feature/Fix/Docs] Brief description of changes
```

### PR Description Template
```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Related Issues
Closes #123
Related to #456

## Changes Made
- [ ] Added authentication middleware
- [ ] Implemented JWT token validation  
- [ ] Added user session management
- [ ] Updated API documentation

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Code review completed

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
```

## Branch Protection Rules

### Master Branch
- Require pull request reviews before merging
- Dismiss stale pull request approvals when new commits are pushed
- Require status checks to pass before merging
- Require branches to be up to date before merging
- Include administrators in restrictions

### Develop Branch  
- Require pull request reviews before merging
- Require status checks to pass before merging
- Require branches to be up to date before merging

## Automated Workflows

### GitHub Actions Integration

```yaml
# .github/workflows/feature-branch.yml
name: Feature Branch CI
on:
  pull_request:
    branches: [ develop ]
    types: [ opened, synchronize, reopened ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21
      - name: Run tests
        run: make test
      - name: Run linting
        run: make lint
```

## Best Practices

### 1. Feature Development
- ✅ Always create feature branches from `develop`
- ✅ Keep feature branches small and focused
- ✅ Regularly merge `develop` into feature branches
- ✅ Use descriptive branch names and commit messages
- ❌ Don't commit directly to `master` or `develop`

### 2. Code Reviews
- ✅ All features must go through PR review
- ✅ At least one approval required before merge
- ✅ Address all review comments before merge
- ✅ Keep PRs reasonably sized (< 500 lines if possible)

### 3. Release Management
- ✅ Use semantic versioning (v1.0.0, v1.1.0, v2.0.0)
- ✅ Create release notes for each version
- ✅ Tag all releases for easy tracking
- ✅ Deploy releases from `master` branch only

### 4. Hotfixes
- ✅ Create from `master` for critical production issues
- ✅ Merge to both `master` and `develop`
- ✅ Create new version tag immediately
- ✅ Deploy hotfix ASAP after testing

This Git flow ensures clean, traceable development while supporting team collaboration and maintaining code quality standards.