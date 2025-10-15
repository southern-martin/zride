#!/bin/bash

# ZRide Git Flow Helper Script
# Usage: ./scripts/git-flow.sh [command] [options]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        error "Not in a git repository"
    fi
}

# Ensure we have the latest changes
update_branch() {
    local branch=$1
    info "Updating $branch branch..."
    git checkout $branch
    git pull origin $branch
}

# Create feature branch
create_feature() {
    local feature_name=$1
    if [ -z "$feature_name" ]; then
        error "Feature name is required. Usage: ./git-flow.sh feature create <feature-name>"
    fi
    
    check_git_repo
    
    local branch_name="feature/$feature_name"
    
    # Check if branch already exists
    if git show-ref --verify --quiet refs/heads/$branch_name; then
        error "Branch $branch_name already exists"
    fi
    
    update_branch "develop"
    
    info "Creating feature branch: $branch_name"
    git checkout -b $branch_name
    git push -u origin $branch_name
    
    success "Feature branch $branch_name created and pushed to remote"
    info "You can now start working on your feature!"
}

# Finish feature (create PR)
finish_feature() {
    check_git_repo
    
    local current_branch=$(git branch --show-current)
    
    if [[ ! $current_branch == feature/* ]]; then
        error "You must be on a feature branch to finish a feature"
    fi
    
    info "Finishing feature: $current_branch"
    
    # Ensure all changes are committed
    if ! git diff --quiet || ! git diff --staged --quiet; then
        error "You have uncommitted changes. Please commit or stash them first."
    fi
    
    # Push latest changes
    git push origin $current_branch
    
    # Update develop and merge
    update_branch "develop"
    
    info "Merging $current_branch into develop"
    git merge $current_branch --no-ff -m "feat: merge $current_branch"
    git push origin develop
    
    # Cleanup
    git branch -d $current_branch
    git push origin --delete $current_branch
    
    success "Feature $current_branch has been merged to develop and cleaned up"
}

# Create release
create_release() {
    local version=$1
    if [ -z "$version" ]; then
        error "Version is required. Usage: ./git-flow.sh release create <version>"
    fi
    
    check_git_repo
    
    local branch_name="release/$version"
    
    update_branch "develop"
    
    info "Creating release branch: $branch_name"
    git checkout -b $branch_name
    
    # Update version in relevant files (you can customize this)
    # echo "$version" > VERSION
    # git add VERSION
    # git commit -m "chore: bump version to $version"
    
    git push -u origin $branch_name
    
    success "Release branch $branch_name created"
    info "Complete your release preparations, then use 'finish_release $version'"
}

# Finish release
finish_release() {
    local version=$1
    if [ -z "$version" ]; then
        error "Version is required. Usage: ./git-flow.sh release finish <version>"
    fi
    
    check_git_repo
    
    local branch_name="release/$version"
    
    if [ "$(git branch --show-current)" != "$branch_name" ]; then
        error "You must be on the $branch_name branch"
    fi
    
    info "Finishing release: $version"
    
    # Merge to master
    git checkout master
    git pull origin master
    git merge $branch_name --no-ff -m "release: version $version"
    git tag -a "v$version" -m "Release version $version"
    git push origin master
    git push origin "v$version"
    
    # Merge back to develop
    git checkout develop
    git pull origin develop
    git merge $branch_name --no-ff -m "chore: merge release $version back to develop"
    git push origin develop
    
    # Cleanup
    git branch -d $branch_name
    git push origin --delete $branch_name
    
    success "Release $version has been completed and tagged as v$version"
}

# Create hotfix
create_hotfix() {
    local fix_name=$1
    if [ -z "$fix_name" ]; then
        error "Hotfix name is required. Usage: ./git-flow.sh hotfix create <fix-name>"
    fi
    
    check_git_repo
    
    local branch_name="hotfix/$fix_name"
    
    update_branch "master"
    
    info "Creating hotfix branch: $branch_name"
    git checkout -b $branch_name
    git push -u origin $branch_name
    
    success "Hotfix branch $branch_name created"
}

# Finish hotfix
finish_hotfix() {
    check_git_repo
    
    local current_branch=$(git branch --show-current)
    
    if [[ ! $current_branch == hotfix/* ]]; then
        error "You must be on a hotfix branch to finish a hotfix"
    fi
    
    info "Finishing hotfix: $current_branch"
    
    # Get version from user
    read -p "Enter hotfix version (e.g., 1.0.1): " version
    
    # Merge to master
    git checkout master
    git pull origin master
    git merge $current_branch --no-ff -m "fix: $current_branch"
    git tag -a "v$version" -m "Hotfix version $version"
    git push origin master
    git push origin "v$version"
    
    # Merge to develop
    git checkout develop
    git pull origin develop
    git merge $current_branch --no-ff -m "fix: merge $current_branch to develop"
    git push origin develop
    
    # Cleanup
    git branch -d $current_branch
    git push origin --delete $current_branch
    
    success "Hotfix $current_branch has been completed and tagged as v$version"
}

# Show current status
show_status() {
    check_git_repo
    
    info "ZRide Git Flow Status"
    echo "===================="
    echo
    
    local current_branch=$(git branch --show-current)
    info "Current branch: $current_branch"
    
    echo
    info "Recent branches:"
    git for-each-ref --sort=-committerdate refs/heads/ --format='%(HEAD) %(color:yellow)%(refname:short)%(color:reset) - %(color:red)%(objectname:short)%(color:reset) - %(contents:subject) - %(authorname) (%(color:green)%(committerdate:relative)%(color:reset))'
    
    echo
    info "Uncommitted changes:"
    if git diff --quiet && git diff --staged --quiet; then
        echo "None"
    else
        git status --short
    fi
}

# Main command handling
case "${1:-}" in
    "feature")
        case "${2:-}" in
            "create")
                create_feature "$3"
                ;;
            "finish")
                finish_feature
                ;;
            *)
                error "Usage: $0 feature [create <name>|finish]"
                ;;
        esac
        ;;
    "release")
        case "${2:-}" in
            "create")
                create_release "$3"
                ;;
            "finish")
                finish_release "$3"
                ;;
            *)
                error "Usage: $0 release [create <version>|finish <version>]"
                ;;
        esac
        ;;
    "hotfix")
        case "${2:-}" in
            "create")
                create_hotfix "$3"
                ;;
            "finish")
                finish_hotfix
                ;;
            *)
                error "Usage: $0 hotfix [create <name>|finish]"
                ;;
        esac
        ;;
    "status")
        show_status
        ;;
    "help"|"--help"|"-h"|"")
        echo "ZRide Git Flow Helper"
        echo "===================="
        echo
        echo "Usage: $0 <command> [options]"
        echo
        echo "Commands:"
        echo "  feature create <name>     Create a new feature branch"
        echo "  feature finish           Finish current feature branch"
        echo "  release create <version>  Create a new release branch"
        echo "  release finish <version>  Finish release branch"
        echo "  hotfix create <name>     Create a new hotfix branch"
        echo "  hotfix finish            Finish current hotfix branch"
        echo "  status                   Show current git flow status"
        echo "  help                     Show this help message"
        echo
        echo "Examples:"
        echo "  $0 feature create user-authentication"
        echo "  $0 feature finish"
        echo "  $0 release create 1.0.0"
        echo "  $0 hotfix create critical-security-fix"
        echo "  $0 status"
        ;;
    *)
        error "Unknown command: $1. Use '$0 help' for usage information."
        ;;
esac