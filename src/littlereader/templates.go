package littlereader

var indexTemplate = `
<!doctype html>
<html>
<head>
<style type="text/css">
body {
	margin: 20px 20px 0 20px;
	background-color: #fafaf8;
	font: 14px 'Lucida Grande', 'Lucida Sans Unicode', Verdana, Helvetica, Arial, Sans-Serif;
	color: #151515;
}
</style>

<script src="http://code.jquery.com/jquery-1.10.1.min.js"></script>
<script>
function hide(s, link) {
	var num = s.split('_')[1];
	$.post('/markAsRead', { href: link });
	$('#' + s).hide();
}

function hideAll(s) {
	$("." + s).each(function(i) {
    $(this).click();
	});
}
</script>
</head>
<body>
<form method="post" action="/add">
Add new subscription: <input type="text" name="url" size=80 />
<input type="submit" value="Add" />
</form>
<br />
`
