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

<div class="row">
  <div class="col-2">
  </div>
  <div class="col-8 bg-light px-0">
      <ol class="breadcrumb">
        <li class="breadcrumb-item"><a href="index">主菜单</a></li>
        <li class="breadcrumb-item"><a href="app">我的应用</a></li>
        <li class="breadcrumb-item active" aria-current="page">应用修改</li>
      </ol>
  </div>
  <div class="col-2">
  </div>
</div>

<div class="row">
<div class="col-2">
</div>
<div class="col-8 bg-light">
    {{if .error}}
    <div class="alert alert-danger">修改失败：{{str2html .error}}</div>
    {{else if .post}}
    <div class="alert alert-success">修改成功</div>
    <br/>
    {{end}}
  <form method="post" action="appmodify" onsubmit="return checkData();">
    <div class="form-group">
      <label for="appname">应用名字：</label>
      <input type="text" class="form-control disable" id="appname" name="appname" value="{{.appname}}" readonly>
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
<div class="col-2">
</div>
</div>

{{template "footer.tpl" .}}