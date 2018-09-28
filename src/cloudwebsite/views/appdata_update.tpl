{{template "header.tpl" .}}

<script type="text/javascript">
    function checkData(){
        if(document.getElementById('nickname').value == ""){
            alert("请输入昵称!");
            return false;
        }
        return true;
    }
</script>

<div class="bg-light">
{{if .post}}
    {{if .error}}
    <div class="alert alert-danger">更新失败：{{str2html .error}}</div>
    {{else if .post}}
    <div class="alert alert-success">更新成功</div>
    <br/>
    {{end}}
{{end}}
  <form method="post" action="update" onsubmit="return checkData();">
    <input type="hidden" name="id" id="id" value="{{.appdata.ID}}" />
    <div class="form-group">
        <label for="appname">应用名称：</label>
        <input type="text" class="form-control disable" id="appname" name="appname" value="{{.appdata.Appname}}" readonly>
    </div>
    <div class="form-group">
        <label for="zonename">分区名：</label>
        <input type="text" class="form-control disable" id="zonename" name="zonename" value="{{.appdata.Zonename}}" readonly>
    </div>
    <div class="form-group">
      <label for="account">账号：</label>
      <input type="text" class="form-control disable" id="account" name="account" value="{{.appdata.Account}}" readonly>
    </div>
    <div class="form-group">
      <label for="nickname">昵称：</label>
      <input type="text" class="form-control" id="nickname" name="nickname" value="{{.appdata.Nickname}}">
    </div>
    <div class="form-group">
      <label for="desc">描述：</label>
      <textarea class="form-control" id="desc" name="desc" rows="3">{{.appdata.Desc}}</textarea>
    </div>
    <div class="form-group">
      <label for="sex">性别：</label>
      <select class="form-control" name="sex" id="sex">
          <option {{if compare .appdata.Sex "男"}}selected{{end}}>男</option>
          <option {{if compare .appdata.Sex "女"}}selected{{end}}>女</option>
      </select>
    </div>
    <div class="form-group">
      <label for="birthday">生日：</label>
      <input type="text" class="form-control" id="birthday" name="birthday" value="{{dateformat .appdata.Birthday "01/02/2006"}}">
    </div>
    <div class="form-group">
      <label for="country">国家：</label>
      <input type="text" class="form-control" id="country" name="country" value="{{.appdata.Country}}">
    </div>
    <button type="submit" class="btn btn-outline-primary btn-lg btn-block mb-3" style="width:100px;">更新</button>
  </form>

</div>

<script type="text/javascript">
  $( function() {
    $( "#birthday" ).datepicker();
    $("#country").countrySelect();
  } );
</script>

{{template "footer.tpl" .}}