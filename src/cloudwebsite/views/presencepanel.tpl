<div class="tab-pane" id="tab_presence" style="max-height:390px;overflow-y:auto;position:absolute;width:245px;">
    <ul id="presencelist" class="products-list product-list-in-box">
        <li id="agree_item" class="item hide">
            <div>friend agree:</div>
            <a href="javascript:void(0)" class="product-title">ID:523455 nickname:wln</a>
            <span class="label label-warning pull-left">2018-06-26T14:29:58.000+08:00</span>
        </li>
        <li id="request_item" class="item hide">
            <div>friend request:</div>
            <a href="javascript:void(0)" class="product-title">ID:123456 nickname:wyq</a>
            <span class="product-description">Hello, friend please</span> 
            <span>                
                <button onclick="agreeFriend($(this).parent().parent().data('data').who);$(this).parent().html('agreed');">add</button>                 
                <button onclick="refuseFriend($(this).parent().parent().data('data').who);$(this).parent().html('refused');">refuse</button>
            </span>
            <span class="label label-warning pull-left">2018-06-26T14:29:58.000+08:00</span>
        </li>
        <li id="request_room_item" class="item hide">
            <div>friend request:</div>
            <a href="javascript:void(0)" class="product-title">ID:123456 nickname:wyq</a>
            <span class="product-description">Hello, friend please</span> 
            <span>                
            <button onclick="reqAgreeRoomJoin($(this).parent().parent().data('data').rid, $(this).parent().parent().data('data').who);$(this).parent().html('agreed');">agree</button>                 
            <button onclick="reqRefuseRoomJoin($(this).parent().parent().data('data').rid, $(this).parent().parent().data('data').who);$(this).parent().html('refused');">refuse</button>
            </span>
            <span class="label label-warning pull-left">2018-06-26T14:29:58.000+08:00</span>
        </li>
    </ul>
</div>

<script>
$( function() {
    $('#tab_presence').hide();
});

function createPresence(data) {
    var item;
    if(data.presencetype == PresenceType.PresenceType_Subscribe){
        if(data.rid){
            item = $("#request_room_item").clone();
            item.attr("id", "request_room_item" + data.who);
        }
        else
        {
            item = $("#request_item").clone();
            item.attr("id", "request_item" + data.who);
        }
    } else {
        item = $("#agree_item").clone();
        item.attr("id", "agree_item" + data.who);
    }

    item.removeClass("hide");
    item.data('data', data);
    var newDate=new Date(parseInt(data.timestamp) * 1000);
    //var html = '<li class="item"> \
    var prestype = "";
    if(data.rid)
        prestype = "room ";
    else
        prestype = "friend "
    
    var reqtype = '';
    if(data.presencetype == PresenceType.PresenceType_Subscribe)
        reqtype += 'request:';
    else if(data.presencetype == PresenceType.PresenceType_Subscribed)
        reqtype += 'agreed:';
    else if(data.presencetype == PresenceType.PresenceType_UnSubscribe)
        if(data.rid)
            reqtype += 'quit:';
        else
            reqtype += 'delete:';
    else if(data.presencetype == PresenceType.PresenceType_UnSubscribed)
        reqtype += 'refuse:';
    else if(data.presencetype == PresenceType.PresenceType_Available)
        reqtype += 'online:';
    else if(data.presencetype == PresenceType.PresenceType_Unavailable)
        reqtype += 'offline:';
    else if(data.presencetype == PresenceType.PresenceType_Invisible)
        reqtype += 'hidden:';
    else
        return null;

    item.find("div").html(prestype + reqtype);
    item.find("a").html('ID:' + data.who + ' nickname:' + data.nickname);
    item.find(".label").html(newDate.format());

    // html += '<a href="javascript:void(0)" class="product-title">';
    // html += 'ID:' + data.who + ' nickname:' + data.nickname;
    // html += '</a>';
    
    // if(data.presencetype == PresenceType.PresenceType_Subscribe){
    //     html += '<span class="product-description">';
    //     html += data.message;
    //     html += '</span> <span>\
    //         <button onclick="agreeFriend(\''+data.who+'\');$(this).parent().html(\'agreed\');">add</button> \
    //         <button onclick="refuseFriend(\''+data.who+'\');$(this).parent().html(\'refused\');">refuse</button></span>';
    // }
        
    // html += '<span class="label label-warning pull-left">' + newDate.format() + '</span></a>';
    // //</li>';

    // var li = $(document.createElement("li"));
    // li.data("user", data);
    // li.addClass("item");
    // li.append(html);

    return item;
}

function addPresence(data) {
    var presence = createPresence(data);
    $("#presencelist").prepend(presence);
}

addPresence({presencetype:0, who:"123456", nickname:"wyq", timestamp:"1529994598", message:"Hello, friend please"});
addPresence({presencetype:1, who:"523455", nickname:"wln", timestamp:"1529994598", message:"Hello, friend please"});
</script>