# Mall Web Application

## Overview
Mall Web is a React 18 application that serves as a frontend for a mall API. It utilizes modern technologies such as React Context API for state management, React Router 6 for routing, Ant Design for UI components, and TailwindCSS for styling. The application is built with Vite and managed with pnpm.

## Features
- User authentication and profile management
- Product listing and detail views
- Cart functionality
- Responsive design with Ant Design and TailwindCSS

## Technologies Used
- React 18
- React Router 6
- Ant Design
- TailwindCSS
- Vite
- pnpm
- TypeScript

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   ```
2. Navigate to the project directory:
   ```bash
   cd mall-web
   ```
3. Install dependencies using pnpm:
   ```bash
   pnpm install
   ```

## Running the Application
To start the development server, run:
```bash
pnpm run dev
```
The application will be available at `http://localhost:3000`.

## Folder Structure
```
mall-web
├── src
│   ├── api
│   ├── components
│   ├── context
│   ├── hooks
│   ├── pages
│   ├── types
│   ├── utils
│   ├── App.tsx
│   ├── main.tsx
│   └── router.tsx
├── .eslintrc.json
├── .gitignore
├── index.html
├── package.json
├── postcss.config.js
├── tailwind.config.js
├── tsconfig.json
└── vite.config.ts
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License.