<div class="tab-pane dragscroll" id="tab_room" style="max-height:390px;overflow-y:auto;position:absolute;width:245px;">
    <div class="box-body" style="">
        <ul class="contacts-list">
        </ul>
    </div>                            
</div>

<script>
function clearRoomList() {
    $("#tab_room .box-body .contacts-list").html("");
}

function addRoom(data) {
    var roomitem = createRoomItem(data);
    $("#tab_room .box-body .contacts-list").append(roomitem);
}

function createRoomItem(data) {
    var li = $(document.createElement("li"));
    li.data("room", data);
    li.dblclick(function(){
        openRoomChatPanel(data);
    });
    li.mouseup(function(e){
        if(e.button===2){
            showRoomMenu(e, data);
            stopPropagation(e);//调用停止冒泡方法,阻止document方法的执行
        }
    });
    //var html = '<li>\
    var html = '<a href="#">\
        <img class="contacts-list-img" src="static/dist/img/user1-128x128.jpg" alt="User Image">\
        <div class="contacts-list-info">\
                <span class="contacts-list-name text-black">';
    html += data.roomname;//     Count Dracula
    html +=     '</span>\
            <span class="contacts-list-msg">';
    html += data.jieshao + '</span>\
        </div>\
        </a>';
    //</li>';
    li.append(html);
    return li;
}
addRoom({rid:"111222", roomname:"adfd", jieshao:"aaaaaaa"});
addRoom({rid:"11122652",roomname:"adfdff", jieshao:"aaaaafaa"});
</script>