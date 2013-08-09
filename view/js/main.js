// Close the dropdown menu if item in menu is clicked
$(document).ready(function closeMenu () {
	$('.top-bar').click(function(){
      $(this).removeClass('expanded');
	});
})
