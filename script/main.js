$(document).ready(function(){
	var nocm = $("#nocm").val()
	
	// $("#nocm").focus(function(){
	// 	var value = $("#nocm").val();
	// 	if (value == ""){
	// 		$("#datapasien").html("Masukkan No. CM");
	// 		//nocm = "";
	// 	} else {
	// 		$("#datapasien").html("No. CM tidak lengkap");
	// 		//nocm = "";
	// 	} 
	// });
	
	$("#navbar").on("keyup", "#nocm", function(){
		
		var value = $("#nocm").val();
		
		
		if (value == ""){
			$("#datapasien").html("Masukkan No. CM");
		} else if (value.length < 8){
			$("#datapasien").html("No. CM tidak lengkap");
		} else {
			$("#nocm").prop("disabled", true);
			nocm = value;
            token = localStorage.getItem("token");

			$.post("getcm", {
				nocm: nocm,
				token: token
			},
			function(data){
				var js = JSON.parse(data)

				if (js.token == "OK"){
					$("#datapasien").html(js.script);
					$("#nocm").prop("disabled", false);
				} else {
					alert("Terjadi kesalahan: " + js.token)
					$("#nocm").prop("disabled", false);
				}			    

			});
        }

	});
	
	$("#navbar").on("click", "#btnsub", function(event){
		event.preventDefault();
		
		var nocm = $("#nocm").val();
		var namapts = $("#namapts").val();
		var diag = $("#diag").val();
		var ats = $("input[name='ats']:checked").val();
		var iki = $("input[name='iki']:checked").val();
		var shift = $("input[name='shift']:checked").val();
		var bagian = $("input[name='bagian']:checked").val();
		var baru = $("#baru").val();
		var dok = $("#email").val();

		switch (true) {
			case nocm === "":
			break;
			case namapts === "":
			break;
			case diag === "":
			console.log("Ini adalah: " + diag);
			break;
			case ats == "":
			console.log("Ini adalah: " + ats);
			break;
			case iki == "":
			console.log("Ini adalah: " + iki);
			break;
			case shift == "":
			console.log("Ini adalah: " + shift);
			break;
			case bagian == "":
			console.log("Ini adalah: " + bagian);
			break;
			default:
			console.log("Default is running");
			// $.post("inputdata",{
			// 	token: localStorage.getItem("token"),
			// 	nocm: nocm,
			// 	namapts: namapts,
			// 	diag: diag,
			// 	ats: ats,
			// 	iki: iki,
			// 	shift: shift,
			// 	bagian: bagian,
			// 	dok: dok,
			// 	baru: baru
			// },
			// function(data){
			// 	$("#nocm").empty();
			// 	var js = JSON.parse(data)
			// 	if (js.token != "OK"){
			// 		alert(js.script)
			// 	}else{
  			// 	    $("tbody").prepend(js.script)
			// 	}


			// });
		}
		alert("Form belum lengkap!")		
	});

	$("#navbar").on("click", "#editbut, #modedit", function(event){
		event.preventDefault();
		token = localStorage.getItem("token")
		var link = $(this).offsetParent().children().first().html();
		$.post("editentri",{
			token: token,
			link: link
		},
		function(data){
			var js = JSON.parse(data)
			$(js.script).modal();
			console.log(js.content.namapts);
			$("body").on('shown.bs.modal', function(){
			var ats = 'input[name="ats"][value=' + js.content.ats + ']';
			var iki = 'input[name="iki"][value=' + js.content.iki + ']';
			var shift = 'input[name="shift"][value=' + js.content.shift + ']';
			var bagian = 'input[name="bagian"][value=' + js.content.bagian + ']';
			if (js.content.bagian == ""){
 		    	$(ats).prop('checked', true);
   				$(iki).prop('checked', true);
   				$(shift).prop('checked', true);	
			
			}else{
 		    	$(ats).prop('checked', true);
   				$(iki).prop('checked', true);
   				$(shift).prop('checked', true);
   				$(bagian).prop('checked', true);

			};
				$('input[name="entri"]').val(js.content.link);
				$('input[name="namapasien"]').val(js.content.namapts);
				$('input[name="diagnosis"]').val(js.content.diag);

			});
			
			
		})

	})

})