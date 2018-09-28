<div style="position:absolute;top:30%;left:50%; transform:translate(-50%, -50%);">  

{{if .post}}
    {{if .error}}
    <div class="alert alert-danger">登录失败：{{str2html .error}}</div>
    <br/>
    {{end}}
{{end}}
<form method="post" action="/main/login" onsubmit="return true;">
  <div class="form-group">
    <label for="account">账号：</label>
    <input type="text" class="form-control" name="account" id="account" placeholder="Account">
  </div>
  <div class="form-group">
    <label for="password1">密码：</label>
    <input type="password" class="form-control" name="password1" id="password1" placeholder="Password" oninput="document.getElementById('password').value = this.value;" onpropertychange="document.getElementById('password').value = this.value;">
    <input type="hidden" name="password" id="password" />
  </div>
  <button type="submit" class="btn btn-outline-primary btn-lg btn-block">登录</button>
</form>

<br/>
<a href="/main/register">点击此处注册</a>

</div>
