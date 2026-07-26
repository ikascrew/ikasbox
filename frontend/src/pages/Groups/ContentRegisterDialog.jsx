import React from "react";

import TextField from '@mui/material/TextField';
import MenuItem from '@mui/material/MenuItem';
import Button from '@mui/material/Button';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

import LoadingButton from "../../components/LoadingButton";
import API from "../../API";

// 生成型コンテンツ(cd / terminal など)を1件登録するダイアログ。
//
// 入力欄はサーバから取得した spec(プラグインの自己申告)から組み立てる。
// このコンポーネントはフィールドの意味を知らず、name / type / label /
// required を見て器を作るだけ。未知の型(fields が空)では生 JSON 入力に
// フォールバックする。
class ContentRegisterDialog extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      view: false,
      types: [],      // [{type, generative, fields}]
      type: "",
      name: "",
      values: {},     // フィールド名 -> 入力値
      rawParams: "",  // 生 JSON 入力の内容
      jsonMode: false, // true なら生 JSON を直接編集する
      error: ""
    };

    this.commitFunc = props.onCommit;
  }

  open = () => {
    this.setState({ view: true, error: "" });
    if (this.state.types.length === 0) {
      this.loadSpec();
    }
  };

  handleClose = () => {
    this.setState({ view: false });
  };

  loadSpec() {
    API.post("/api/v1/contents/spec", {}).then((res) => {
      const types = res.data.types || [];
      // 既定は最初の生成型(実体ファイルが要らないもの)。
      // file / img はディレクトリ import が主経路なので初期選択にしない
      const first = types.find((t) => t.generative) || types[0];
      this.setState({ types: types });
      if (first) {
        this.selectType(first.type, types);
      }
    }).catch((err) => {
      this.setState({ error: err.message });
    });
  }

  // selectType は型を切り替え、その型の既定値で入力値を作り直す。
  // フォーム定義を持たない型(Spec が nil = 未知のプラグイン)では
  // 生 JSON 編集へ自動的に切り替える
  selectType(type, types) {
    const list = types || this.state.types;
    const entry = list.find((t) => t.type === type);

    const values = {};
    const hasFields = !!(entry && entry.fields && entry.fields.length > 0);

    if (hasFields) {
      entry.fields.forEach((f) => {
        values[f.name] = f.default || "";
      });
    }

    this.setState({
      type: type,
      values: values,
      jsonMode: !hasFields,
      rawParams: hasFields ? "" : "{}",
      error: ""
    });
  }

  // handleToggleJson はフォームと生 JSON を双方向に引き継いで切り替える。
  // JSON へ移るときは今のフォーム内容を、フォームへ戻るときは JSON の
  // 内容を(読めれば)入力値に反映する
  handleToggleJson = (event) => {

    const on = event.target.checked;

    if (on) {
      let text = "{}";
      try {
        text = JSON.stringify(JSON.parse(this.buildParams()), null, 2);
      } catch {
        // 組み立てに失敗したら空のオブジェクトから始める
      }
      this.setState({ jsonMode: true, rawParams: text, error: "" });
      return;
    }

    const entry = this.currentEntry();
    const values = { ...this.state.values };
    try {
      const obj = JSON.parse(this.state.rawParams);
      if (entry && entry.fields) {
        entry.fields.forEach((f) => {
          if (obj[f.name] !== undefined) {
            values[f.name] = String(obj[f.name]);
          }
        });
      }
    } catch {
      // JSON が壊れている場合はフォーム側の値をそのまま残す
    }

    this.setState({ jsonMode: false, values: values, error: "" });
  };

  handleChangeType = (event) => {
    this.selectType(event.target.value);
  };

  handleChangeName = (event) => {
    this.setState({ name: event.target.value });
  };

  handleChangeField = (name) => (event) => {
    this.setState({
      values: { ...this.state.values, [name]: event.target.value }
    });
  };

  handleChangeRaw = (event) => {
    this.setState({ rawParams: event.target.value });
  };

  currentEntry() {
    return this.state.types.find((t) => t.type === this.state.type);
  }

  // buildParams は入力値を params の JSON 文字列に組み立てる。
  // 空欄のフィールドは送らない(プラグイン側の既定に任せる)
  buildParams() {
    const entry = this.currentEntry();
    if (this.state.jsonMode || !entry || !entry.fields || entry.fields.length === 0) {
      return this.state.rawParams;
    }

    const obj = {};
    entry.fields.forEach((f) => {
      const v = this.state.values[f.name];
      if (v !== undefined && v !== "") {
        obj[f.name] = v;
      }
    });
    return JSON.stringify(obj);
  }

  handleRegister = () => {

    const entry = this.currentEntry();

    if (this.state.name === "") {
      this.setState({ error: "Name is required." });
      return Promise.reject("error");
    }

    if (this.state.jsonMode) {
      // 壊れた JSON はサーバへ行く前に弾く(サーバ側は plugin 生成で
      // 落ちるため、ここで出す方がメッセージが分かりやすい)
      try {
        JSON.parse(this.state.rawParams);
      } catch (e) {
        this.setState({ error: "Params is not valid JSON: " + e.message });
        return Promise.reject("error");
      }
    } else if (entry && entry.fields) {
      // 必須フィールドの未入力はサーバへ行く前に弾く
      const missing = entry.fields.find(
        (f) => f.required && !this.state.values[f.name]
      );
      if (missing) {
        this.setState({ error: (missing.label || missing.name) + " is required." });
        return Promise.reject("error");
      }
    }

    const args = {
      groupId: Number(this.props.groupId),
      name: this.state.name,
      type: this.state.type,
      params: this.buildParams()
    };

    return API.post("/api/v1/contents/register", args).then(() => {
      if (this.commitFunc !== undefined) {
        this.commitFunc();
      }
      this.setState({ name: "", error: "" });
      this.handleClose();
    }).catch((err) => {
      // params が不正な場合はサーバ側(プラグイン生成)で弾かれる
      this.setState({ error: err.message });
      throw err;
    });
  };

  // renderField は param.Type に対応する入力欄を返す。
  // 対応表を増やす場合は plugin 側の param.Type も併せて変更すること
  renderField(f) {

    const value = this.state.values[f.name] || "";

    switch (f.type) {
      case "multiline":
        return (
          <TextField
            key={f.name}
            margin="dense"
            label={f.label || f.name}
            variant="standard"
            value={value}
            onChange={this.handleChangeField(f.name)}
            required={f.required}
            multiline
            minRows={3}
            fullWidth
          />
        );
      case "datetime":
        return (
          <TextField
            key={f.name}
            margin="dense"
            type="datetime-local"
            label={f.label || f.name}
            variant="standard"
            value={value}
            onChange={this.handleChangeField(f.name)}
            required={f.required}
            InputLabelProps={{ shrink: true }}
            fullWidth
          />
        );
      default:
        return (
          <TextField
            key={f.name}
            margin="dense"
            type="text"
            label={f.label || f.name}
            variant="standard"
            value={value}
            onChange={this.handleChangeField(f.name)}
            required={f.required}
            fullWidth
          />
        );
    }
  }

  render() {

    const entry = this.currentEntry();
    const fields = (entry && entry.fields) || [];

    return (<>
      <Dialog open={this.state.view} fullWidth maxWidth="sm">
        <DialogTitle>Add Content</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Register a single content by parameters. The form below is
            defined by the selected plugin.
          </DialogContentText>

          <TextField
            select
            margin="dense"
            label="Type"
            variant="standard"
            value={this.state.type}
            onChange={this.handleChangeType}
            fullWidth
          >
            {this.state.types.map((t) => (
              <MenuItem key={t.type} value={t.type}>
                {t.type}{t.generative ? "" : " (file required)"}
              </MenuItem>
            ))}
          </TextField>

          <TextField
            autoFocus margin="dense"
            type="text" label="Name"
            variant="standard"
            value={this.state.name}
            onChange={this.handleChangeName}
            required
            fullWidth
          />

          {/* フォーム定義を持たない型では JSON 編集しか選べない */}
          <FormControlLabel
            sx={{ marginTop: 1 }}
            control={
              <Switch
                checked={this.state.jsonMode}
                onChange={this.handleToggleJson}
                disabled={fields.length === 0}
              />
            }
            label="Edit params as JSON"
          />

          {this.state.jsonMode
            ? (
              <TextField
                margin="dense"
                label="Params (JSON)"
                variant="standard"
                value={this.state.rawParams}
                onChange={this.handleChangeRaw}
                multiline
                minRows={4}
                fullWidth
                helperText={
                  fields.length === 0
                    ? "This type has no form definition. Write the params JSON directly."
                    : "The form fields are ignored while this is on."
                }
              />
            )
            : fields.map((f) => this.renderField(f))}

          {this.state.error !== "" && (
            <DialogContentText sx={{ color: "error.main", marginTop: 1 }}>
              {this.state.error}
            </DialogContentText>
          )}

        </DialogContent>

        <DialogActions>
          <Button onClick={this.handleClose}>Cancel</Button>
          <LoadingButton onClick={this.handleRegister}>Register</LoadingButton>
        </DialogActions>
      </Dialog>
    </>);
  }
}

export default ContentRegisterDialog;
