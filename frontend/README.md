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

`vite.config.js` proxies `/api`, `/thumb`, and `/content/media` to `http://localhost:5555` (the Go server) and listens on port 3000. Only the `/content/media` prefix is proxied, not the bare `/content` — that would also swallow the frontend's own `/contents` SPA route.

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

- グループ／プロジェクトの削除、プロジェクト名の変更（未実装）
- ディレクトリインポートの進捗表示（現状は非同期で投げっぱなし）
