# installation

```
$ npm install
```

# structure

- frontend
    - index.html
    - src
        - App.jsx
    - vite.config.js

Built with Vite + `@vitejs/plugin-react`. Files containing JSX use the `.jsx` extension (`App.jsx`, `pages/*.jsx`, `components/*.jsx`); plain modules with no JSX (`API.js`, `Util.js`, `Paging.js`) stay `.js`.

# client server

```
$ npm run dev
```

http://localhost:3000/

`vite.config.js` proxies `/api` and `/thumb` to `http://localhost:5555` (the Go server) and listens on port 3000.

# build

```
$ npm run build
```

Outputs to `frontend/dist` (gitignored). `npm run preview` serves that build locally.

# Design

## Material-Design(Lite)

```
$ npm install @mui/material @emotion/react @emotion/styled
```

```
<link
  rel="stylesheet"
  href="https://fonts.googleapis.com/icon?family=Material+Icons"
/>
```

# Router

```
$ npm install react-router-dom
```

# API Access

```
$ npm install --save axios
```

# Session(Cookie)

not implemented.

# Test

not implemented.


# issue

- コンテンツ一覧
- プロジェクト作成
- プロジェクトグループ追加
