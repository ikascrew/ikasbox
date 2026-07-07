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

Outputs to `../handler/internal/_assets/spa` (gitignored except a `.gitkeep` placeholder), which the Go server `go:embed`s and serves directly — `go run main.go start` alone is enough to see the production build, no separate static server needed. Run this at least once on a fresh checkout before building the Go side, since the embed directive needs at least one file to exist there. `npm run preview` serves the same build locally via Vite instead.

# Design

## Material UI (MUI)

The whole app is themed through a single `createTheme()` in `src/theme.js`, applied via `ThemeProvider` + `CssBaseline` in `App.jsx`. There is no separate custom layout CSS — `pages/Layout.jsx` (AppBar + permanent Drawer) and everything else style through the theme and MUI's `sx` prop. The old MDL (Material Design Lite) based legacy pages have been removed entirely.

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

- グループ／プロジェクト名の変更（未実装。削除は実装済み）
- ディレクトリインポートの進捗表示（現状は非同期で投げっぱなし）
