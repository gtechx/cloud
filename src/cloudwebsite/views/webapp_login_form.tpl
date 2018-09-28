<div style="position:absolute;top:40%;left:50%; transform:translate(-50%, -50%);">  

  <div class="alert alert-danger hide" id="error">dfdddd</div>
  <br/>

  <div id="loginpanel" class="box box-widget">
  <div class="overlay hide">
    <i class="fa fa-refresh fa-spin"></i>
  </div>
<div class="login-box">
<div class="login-box-body">
<form class="animated" onsubmit="return false;">
  <div class="form-group">
    <label for="appname">选择应用：</label>
    <select class="form-control" onchange="onAppnameChange(this);" name="appname" id="appname">
        {{range $index, $elem := .applist}}
        <option>{{$elem.Appname}}</option>
        {{end}}
    </select>
  </div>
  <div class="form-group">
    <label for="zonename">选择分区：</label>
    <select class="form-control" name="zonename" id="zonename">
        <option></option>
    </select>
  </div>
  <div class="form-group">
    <label for="platform">选择平台：</label>
    <select class="form-control" name="platform" id="platform">
        <option>web</option>
        <option>web1</option>
        <option>web2</option>
        <option>web2</option>
    </select>
  </div>
  <div class="form-group">
    <label for="account">账号：</label>
    <input type="text" class="form-control" name="account" id="account" placeholder="Account" value="{{if .appaccount}}{{.appaccount}}{{else}}wyq{{end}}" />
  </div>
  <div class="form-group">
    <label for="password1">密码：</label>
    <input type="password" class="form-control" name="password1" id="password1" placeholder="Password" oninput="document.getElementById('password').value = this.value;" onpropertychange="document.getElementById('password').value = this.value;" value="123" />
    <input type="hidden" name="password" id="password" value="123" />
  </div>
  <button onclick="dologin(); return false;" class="btn btn-outline-primary btn-lg btn-block">登录</button>
</form>
</div>
</div>

  <script type="text/javascript">
    $(function () {
      onAppnameChange(document.getElementsByName("appname")[0]);
    });

    function onAppnameChange(obj) {
      if (obj.selectedIndex == -1)
        return;
      var opt = obj.options[obj.selectedIndex];

      $.post("/webapp/zonelist", { 'appname': $('#appname').val() },
        function (data) {
          var jsondata = JSON.parse(data);
          var liststr = '';
          var count = jsondata["total"];
          var html = "";
          for (i in jsondata["rows"]) {
            var row = jsondata["rows"][i];
            html += '<option>' + row.zonename + '</option>';
          }
          $('#zonename').html(html);
        });
    };

    function dologin() {
      setPlatform($('#platform').val());
      login($('#account').val(), $('#password').val(), $('#appname').val(), $('#zonename').val());
    }
  </script>
  </div>

  <div id="idselect" class="hide">
    <div id="idlist" class="box-body">
    <button type="button" class="btn btn-default btn-block">.btn-block</button>
    <button type="button" class="btn btn-default btn-block">.btn-block</button>
    </div>
    <button onclick="quitChat();" class="btn btn-outline-primary btn-lg btn-block">退出</button>
  </div>
</div>



