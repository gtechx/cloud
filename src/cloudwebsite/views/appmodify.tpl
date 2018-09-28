{{template "header.tpl" .}}

<script type="text/javascript">
    function checkData(){
        if(document.getElementById('appname').value == ""){
            alert("请输入应用名字!");
            return false;
        }
        return true;
    }
</script>

<div class="bg-light">
    {{if .error}}
    <div class="alert alert-danger">修改失败：{{str2html .error}}</div>
    {{else if .post}}
    <div class="alert alert-success">修改成功</div>
    <br/>
    {{end}}
  <form method="post" action="update" onsubmit="return checkData();">
    <div class="form-group">
      <label for="appname">应用名字：</label>
      <input type="text" class="form-control disable" id="appname" name="appname" value="{{.appname}}" {{if not .isadmin}}readonly{{end}}>
    </div>
    <div class="form-group">
      <label for="desc">应用介绍：</label>
      <textarea class="form-control" id="desc" name="desc" rows="3">{{.desc}}</textarea>
    </div>
    <div class="form-group">
      <label for="share">共享数据应用名字：</label>
      <input type="text" class="form-control" id="share" name="share"  value="{{.share}}">
    </div>    
    <button type="submit" class="btn btn-outline-primary btn-lg btn-block mb-3" style="width:100px;">修改</button>
  </form>

</div>

{{template "footer.tpl" .}}