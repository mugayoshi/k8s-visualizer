# K8s Visualizer Frontend

This is the frontend for the K8s Visualizer application, built with SvelteKit and Vite.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   [Node.js](https://nodejs.org/) (v18 or newer)
*   [Yarn](https://classic.yarnpkg.com/)

### Installation

1.  Clone the repo:
    ```sh
    git clone https://github.com/mugayoshi/k8s-visualizer.git
    ```
2.  Navigate to the frontend directory:
    ```sh
    cd k8s-visualizer/frontend
    ```
3.  Install the dependencies:
    ```sh
    yarn install
    ```
4.  Enable the pre-commit hook:
    ```sh
    git config core.hooksPath frontend/.husky
    ```
    This will ensure that tests are run before each commit.

## Available Scripts

In the project directory, you can run:

### `yarn dev`

Runs the app in development mode.
Open [http://localhost:5173](http://localhost:5173) to view it in your browser.

The page will reload when you make changes.
You may also see any lint errors in the console.

### `yarn build`

Builds the app for production to the `.svelte-kit` directory.
It correctly bundles the app in production mode and optimizes the build for the best performance.

### `yarn preview`

Runs a local server to preview the production build.

### `yarn test`

Runs all tests once. This is also used by the pre-commit hook.

### `yarn test:unit`

Launches the test runner in interactive watch mode.

### `yarn lint`

Lints the project files for any code quality issues.

### `yarn format`

Formats all project files with Prettier.

### `yarn storybook`

Runs Storybook, a tool for UI development.
Open [http://localhost:6006](http://localhost:6006) to view it in your browser.
