var tradeNo;
var i = 0;
$(function () {
    tradeNo = LinkUrl.Request("out_trade_no") || "";
    if (tradeNo != "") {
        queryOrder()
    } else {
        layer.msg("无效订单编号", {icon: 1}, function () {
            var backUrl = Cookies.get('backUrl') || false;
            if (backUrl) {
                location.href = backUrl
            } else {
                location.href = "/"
            }
        })
    }
});

function queryOrder() {
    $.ajax({
        load_tip: false,
        url: "/order/query",
        data: {tradeNo: tradeNo},
        success: function (d) {
            if (d.ret == 200) {
                if (d.data.type == "会员套餐") {
                    if (d.data.member_type_id == 1) {
                        location.href = "/app/page/notify/success.html"
                    } else {
                        location.href = "/app/page/notify/success-company.html"
                    }
                } else {
                    layer.msg("支付成功", {icon: 1}, function () {
                        var backUrl = Cookies.get('backUrl') || false;
                        if (backUrl) {
                            location.href = backUrl
                        } else {
                            location.href = "http://www.tbqbz.com/"
                        }
                    })
                }
            }
            else if (d.ret == 404) {
                if (i < 900) {
                    setTimeout(function () {
                        i++
                        queryOrder()
                    }, 1000)
                } else {
                    //15分钟后不再查询
                }
            }
        }
    });
}