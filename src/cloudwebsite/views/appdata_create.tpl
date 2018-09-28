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
    <div class="alert alert-danger">创建失败：{{str2html .error}}</div>
    <br/>
    {{end}}
{{end}}
  <form method="post" action="create" onsubmit="return checkData();">
    <div class="form-group">
        <label for="appname">应用名称：</label>
        <select class="form-control" onchange="onAppnameChange(this);" style="width:100px" name="appname" id="appname">
            <option></option>
            {{range $index, $elem := .applist}}
            <option>{{$elem.Appname}}</option>
            {{end}}
        </select>
    </div>
    <div class="form-group">
        <label for="zonename">分区名：</label>
        <select class="form-control" style="width:100px" name="zonename" id="zonename">
            <option></option>
        </select>
    </div>
    <div class="form-group">
      <label for="account">账号：</label>
      <input type="text" class="form-control" id="account" name="account">
    </div>
    <div class="form-group">
      <label for="nickname">昵称：</label>
      <input type="text" class="form-control" id="nickname" name="nickname">
    </div>
    <div class="form-group">
      <label for="desc">描述：</label>
      <textarea class="form-control" id="desc" name="desc" rows="3"></textarea>
    </div>
    <div class="form-group">
      <label for="sex">性别：</label>
      <select class="form-control" name="sex" id="sex">
          <option>男</option>
          <option>女</option>
      </select>
    </div>
    <div class="form-group">
      <label for="birthday">生日：</label>
      <input type="text" class="form-control" id="birthday" name="birthday">
    </div>
    <div class="form-group">
      <label for="country">国家：</label>
      <input type="text" class="form-control" id="country" name="country">
    </div>
    <button type="submit" class="btn btn-outline-primary btn-lg btn-block mb-3" style="width:100px;">创建</button>
  </form>

</div>

<script type="text/javascript">
  $( function() {
    $( "#birthday" ).datepicker();
    $("#country").countrySelect();
  } );

  function onAppnameChange(obj) {
      var opt = obj.options[obj.selectedIndex];
      console.info("text:"+opt.text);
      console.info("value:"+opt.value);

        $.post("zonelist", { 'account': $('#account').val(), 'appname': $('#appname').val()},
        function(data) {
        $('#zonename').bootstrapTable('refresh');
        var jsondata = JSON.parse(data);
        console.info(jsondata["rows"]);
        var liststr = '';
        var count = jsondata["total"];
        var html = $('#zonename').html();
        for (i in jsondata["rows"])
        {
            var row = jsondata["rows"][i];
            html += '<option>' + row.zonename + '</option>';
        }
        $('#zonename').html(html);
        });
  }
</script>

{{template "footer.tpl" .}}