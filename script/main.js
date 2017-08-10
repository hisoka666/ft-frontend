
$(document).ready(function(){
	var nocm = $("#nocm").val()
	
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
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Diagnosis belum diisi" +
		                        "</div>");
			// alert("Diagnosis harus diisi");
			break;
			case ats == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"ATS belum diisi" +
		                        "</div>");
			// alert("ATS harus diisi");
			break;
			case iki == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Golongan IKI belum diisi" +
		                        "</div>");
			// alert("Golongan IKI harus diisi");
			break;
			case shift == null :
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Shift belum diisi" +
		                        "</div>");
			// alert("Shift harus diisi");
			break;
			case bagian == null:
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		                            "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
									"Bagian belum diisi" +
		                        "</div>");
			// alert("Bagian harus diisi");
			break;
			default:
			$.post("inputdata",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				dok: dok,
				baru: baru
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					alert(js.script);
				}else{
					$("tbody").prepend(js.script);
					refreshNumber();
					$("#nocm").val('');
					$("#datapasien").html('');
					$("tbody tr:eq(101)").remove();
					popModalWarning("Sukses", "Berhasil menambahkan data", "")
					
				}
			});
		}
	});

	$("#navbar").on("click", "#editbut", function(event){
		event.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		token = localStorage.getItem("token");
			$.post("editentri",{
				token: token,
				link: link
			},
			function(data){
				var js = JSON.parse(data);
				$("#mymodal").html(js.script);
				var ats = 'input[name="modats"][value=' + js.content.ats + ']';
				var iki = 'input[name="modiki"][value=' + js.content.iki + ']';
				var shift = 'input[name="modshift"][value=' + js.content.shift + ']';
				var bagian = 'input[name="modbagian"][value=' + js.content.bagian + ']';
			
				$('input[name="entri"]').val(js.content.link);
				$('input[name="namapasien"]').val(js.content.nama);
				$('input[name="diagnosis"]').val(js.content.diag);
				$('input[name="urutan"]').val(indexrow)
			
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
				$("#mymodal").modal();
			});
	});

$('body').on('hide.bs.modal', "#mymodal" , function(){
	console.log("Removal mymodal works");
	$(this).removeData('bs.modal').find(".modal-content").empty();
	});
$('body').on('hide.bs.modal', "#modwarning" , function(){
	console.log("removal modwarning works");
	$(this).removeData('bs.modal');
	});

$("body").on("click", "#confirmedit", function(e){
		e.preventDefault();
		
		var link = $("#modentri").val();
		var namapts = $('input[name="namapasien"]').val();
		var diag = $('input[name="diagnosis"]').val();
		var ats = $("input[name='modats']:checked").val();
		var bagian = $("input[name='modbagian']:checked").val();
		var iki = $("input[name='modiki']:checked").val();
		var shift = $("input[name='modshift']:checked").val();
		var indexrow = $('input[name="urutan"]').val();
		var urutan = "tbody.dafpts tr:eq(" + $('input[name="urutan"]').val() + ")";

		
		switch (true) {
			case namapts === "":
			popModalWarning();
			break;
			case diag === "":
			popModalWarning();
			break;
			case ats == null :
			break;
			case iki == null :
			break;
			case shift == null :
			break;
			case bagian == null:
			break;
			default:
			$.post("confedit",{
				token: localStorage.getItem("token"),
				nocm: nocm,
				namapts: namapts,
				diag: diag,
				ats: ats,
				iki: iki,
				shift: shift,
				bagian: bagian,
				link: link
			},
			function(data){
				var js = JSON.parse(data)
				if (js.token != "OK"){
					popModalWarning("Kesalahan Pada Server", "Telah terjadi kesalahan pada server. Mohon ulangi proses sebelumnya");
				}else{
					$(urutan).replaceWith(js.script);
					console.log("Editing jalan")
					refreshNumber();
					$("#mymodal").modal('hide');
					popModalWarning("Edit Entri", "Berhasil mengubah entri", "");
					
				}
			});
		}

	});

var refreshNumber = function(){
	
$("tr").find(".nourut").each(function(index, elem){
		num = index + 1;
		$(elem).html(num)
	})
}

var popModalWarning = function(title, msg, prop){
		$(".modal-title").html(title);
		$("#message").html(msg);
		if (prop == ""){
			$("#extrabut").hide();
			$("#modwarning").modal();
		} else {
			$("#extrabut").html(prop);
			$("#extrabut").show();
			$("#modwarning").modal();
		}			
};

$("#navbar").on("click", "#delbut", function(e){
		e.preventDefault();
		var link = $(this).offsetParent().children().first().html();
		var indexrow = $(this).closest("tr").index();
		var urutan = "tbody tr:eq(" + indexrow + ")";
		var token = localStorage.getItem("token");
		popModalWarning("Hapus Entri", "Yakin ingin menghapus entri ini?", "Hapus");
		$("body").one("click", "#extrabut", function(){
			
			$.post("confdel", {
				link: link,
				token: token
			}, function(data){
				var js = JSON.parse(data);
				if (js.token == "OK"){
					$(urutan).remove();
					console.log("Removed : "+ urutan);
					refreshNumber();
					$("#modwarning").modal('hide');
				}else{
					alert("Terjadi kesalahan");
				}
			});
		});
	});






$("#navbar").on("click", "#homebutton", function(e){
		token = localStorage.getItem("token");
		email = $("#email").val();
		e.preventDefault();
		$.post("/firstentries", {
			token: token,
			email: email
		},function(data){
			var js = JSON.parse(data);
			// console.log(js.script);
			if (js.token == "OK") {
				$("#inputnocm").show();
				$("#tabeliki").hide();
				$(".diagram").hide();
				$("#tabelutama").html(js.script);
				// removeModal("#modwarning")
			}else{
				popModalWarning("Peringatan", "Terjadi kesalahan pada server. Hubungi admin", "");
			}
			

		})

	});

$("#navbar").on("click", "#tglbut", function(e){
	token = localStorage.getItem("token");
	var link = $(this).offsetParent().children().first().html();
	e.preventDefault();
	console.log(link);

	$.post("edittgl", {
		token: token,
		link: link
	}, function(data){
		var js = JSON.parse(data);
		if (js.token != "OK") {
			$("#mymodal").html(js.modal);
		}else{
			$("#mymodal").html(js.script);
			$("#inputdate").datepicker({
				dateFormat: "dd-mm-yy",
				altField: "#datesend",
				altFormat: "yy-mm-dd"
			})
		}
		$("#mymodal").modal();
	})
});

$("body").on("click", "#ubahtglbut", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tanggal = $("#datesend").val();
	var link = $("#linkubahtgl").val();
	console.log("tanggal adalah : " + tanggal);
	// var tglasli = $(this).find("#tglasli").html();
	$.post("confedittgl", {
		token: token,
		link: link,
		tanggal: tanggal
	}, function(data){
		var js = JSON.parse(data);
		// console.log("Isi dari script adalah: " + data.script)
		$("#mymodal").modal('hide');
		if (js.token != "OK") {			
			popModalWarning("Peringatan", "Terjadi kesalahan pada server. Hubungi admin","");
		} else {
			$("#tabelutama").html(js.script);
			popModalWarning("Sukses", "Berhasil mengubah tanggal","");
		}

	})
	// console.log(jQuery.type(tanggal));
	// $("#mymodal").modal('hide');
});

$("#navbar").on("click", "#makeresep", function(e){
	e.preventDefault();

	$.get("getptspage")
	.done(function(data){
		var js = JSON.parse(data);
		$("#mymodal").html(js.script);
		$("#mymodal").modal();
		$("#tgllahir").datepicker({
			dateFormat:"dd-mm-yy",
			changeMonth: true,
			changeYear: true,
			yearRange: "1900:2035",
			onSelect: function(value, ui){
				var today = new Date();
				// console.log("Tahun ini adalah: " + today.getFullYear());
				// console.log("Tahun yang dipilih adalah : " + ui.selectedYear)
				var umur = today.getFullYear() - ui.selectedYear;
				// console.log("Umur adalah: " + umur);
				$("#umur").val(umur);
			}
		});
	});

	$("body").on("click", "#rspnextbut", function(e){
		e.preventDefault();
		
		var nama = $("#namapts").val();
		var diag = $("#diag").val();
		var umur = $("#umur").val();
		var almt = $("#almt").val();
		var bb = $("#bb").val();
		var alergi = $("#alergi").val();


		$("#mymodal").modal();
		$.get("getprespage")
		.done(function(resep){
			var jso = JSON.parse(resep);
			$("#mymodal").html(jso.script);
			$("#rspnamapts").html(nama);
			$("#rspdiag").html(diag);
			$("#rspumur").html(umur);
			$("#rspbb").html(bb);
			$("#rspalmt").html(almt);
			$("#rspalergi").html(alergi);
			$("#mymodal").modal();
		})
		
		
	})
});

$("#navbar").on("click", "#inputobat", function(e){
	e.preventDefault();
	$.get("getinputobat")
	.done(function(data){
		var js = JSON.parse(data);
		if (js.token != "OK"){
			popModalWarning("Peringatan", "Gagal memperoleh template", "")
		}else{
			$("#mymodal").html(js.script)
			$("#mymodal").modal();
		}

	})
})

$("body").on("click", ".btn.tablet.tambah", function(e){
	tambahElement("tablet", this)
	// $("div.tablet.col-sm-9").last().clone().appendTo(".form-group.tablet");
	// $("div.tablet.col-sm-3").last().clone().appendTo(".form-group.tablet");
	// $("input.form-control.tablet").each(function(){
	// 	$(this).prop("disabled",true);
	// });
	// $("input.form-control.tablet").last().val("").prop("disabled", false);
	// $(this).html("Hapus").removeClass("tambah").addClass("hapus");
});

$("body").on("click", ".btn.tablet.hapus", function(e){
	e.preventDefault();
	hapusElement("tablet", this);
	// var index = $(".tablet.col-sm-3").index($(this).parent());
	// console.log("index adalah : " + index);
	// $(".tablet.col-sm-9").eq(index).remove();
	// $(".tablet.col-sm-3").eq(index).remove();
});

$("body").on("click", ".btn.sirop.tambah", function(e){
	e.preventDefault();
	tambahElement("sirop", this);
});

$("body").on("click", ".btn.sirop.hapus", function(e){
	e.preventDefault();
	hapusElement("sirop",this)
});

$("body").on("click", ".btn.drop.tambah", function(e){
	e.preventDefault();
	tambahElement("drop", this);
});

$("body").on("click", ".btn.drop.hapus", function(e){
	e.preventDefault();
	hapusElement("drop",this)
});
var tambahElement = function(selector, elem){
	$("div."+ selector + ".col-sm-9").last().clone().appendTo(".form-group."+selector);
	$("div."+ selector + ".col-sm-3").last().clone().appendTo(".form-group."+selector);
	$("input.form-control."+selector).each(function(index, elem){
		// $(this).prop("disabled",true);
		$(elem).attr("name", index);
	});
	$("input.form-control."+selector).last().val("");
	$(elem).html("Hapus").removeClass("tambah").addClass("hapus");

}

var hapusElement = function(sel, elem){
	var index = $("."+sel+".col-sm-3").index($(elem).parent());
	$("."+sel+".col-sm-9").eq(index).remove();
	$("."+sel+".col-sm-3").eq(index).remove();
}


$("body").on("click", "#savdrug", function(e){
	e.preventDefault();
	if ($("#mrkdgng").val() == ""||
		$("#kand").val() == ""||
		$("#mindose").val() == undefined || 
		$("#maxdose").val() == undefined || 
		$("input.tablet").serializeArray().length == 0 || 
		$("input.sirop").serializeArray().length == 0 ||
		$("input.drop").serializeArray().length == 0 ||
		$("#rekom").val() == "" ) {
			$("#alertmsg").html("<div class=\"alert alert-danger alert-dismissable\"\>" +
		    	                "<a href=\"#\" class=\"close\" data-dismiss=\"alert\" aria-label=\"close\">&times;</a\>" +
								"Form tidak boleh kosong" +
		                        "</div>");
		}else{
			$.post("inputobat", {
				merk: $("#mrkdgng").val(),
				kand: $("#kand").val(),
				mindose: $("#mindose").val(),
				maxdose: $("#maxdose").val(),
				tab : $("input.tablet").serializeArray(),
				syr : $("input.sirop").serializeArray(),
				drop: $("input.drop").serializeArray(),
				rekom : $("#rekom").val(),
				doc : $("#email").val()
			}, function(data){
				var js = JSON.parse(data);
			})
		}
	// var merk = $("#mrkdgng").val();
	// var kandungan = $("#kand").val();
	// var mindose = $("#mindose").val();
	// var maxdose = $("#maxdose").val();
	// var tablet = $("inputtablet").val();
	// console.log(string($("input").serializeArray()));
	var tablet = $("input.tablet").serializeArray();
	var sirop = $("input.sirop").serializeArray();
	var drop = $("input.drop").serializeArray();
	var tab = JSON.stringify(tablet);
	var syr = JSON.stringify(sirop);
	var gtt = JSON.stringify(drop);
	console.log(tab);
	console.log(syr);
	console.log(gtt);	
	// jQuery.each(obat, function(k, v){
	// 	console.log(k + " adalah " + v);
	// })

});

$("#navbar").on("click", "#bulanini", function(e){
	e.preventDefault();
	var now = new Date();
	var dateone = new Date(now.getFullYear(),now.getMonth(),1,8,0,0);
	var token = localStorage.getItem("token");
	if (now > dateone){
		$.post("getmonthly", {
			token: token,
			month: now.getMonth() + 1,
			year: now.getFullYear(),
			email: $("#email").val()
		}, function(data){
			var js = JSON.parse(data);
			pieChart(js.data, "")
			$("#inputnocm").hide();
			$(".diagram").show();
			$("#tabeliki").html(js.modal).show();
			getSum();
			$("#tabelutama").html(js.script);
		})
	}else{
		var blnlalu = new Date(now.getFullYear(), now.getMonth() - 1, 1)
		$.post("getmonthly", {
			token: token,
			month: blnlalu.getMonth() + 1,
			year: blnlalu.getFullYear(),
			email: $("#email").val()
		}, function(data){
			var js = JSON.parse(data);
			pieChart(js.data, "")
			$(".diagram").show();
			$("#inputnocm").hide();
			$("#tabeliki").html(js.modal).show();
			getSum();
			$("#tabelutama").html(js.script);
		})
	}


});

var getSum = function(){
	var sum = 0;
	var iki11 = 0;
	var iki12 = 0;
	for (i=0;i<16;i++){
		iki11 = iki11 + parseInt($("td.jml").eq(i).html())
	}
	for (i=16;i<32;i++){
		iki12 = iki12 + parseInt($("td.jml").eq(i).html())
	}
	var iki21 = iki11;
	var iki22 = iki12;
	for (i=32;i<47;i++){
		iki21 = iki21 + parseInt($("td.jml").eq(i).html())
	}
	for (i=47;i<62;i++){
		iki22 = iki22 + parseInt($("td.jml").eq(i).html())
	}
	
	var tot1 = iki21 * 0.0032;
	var tot2 = iki22 * 0.01;
	var totfinal = tot1 + tot2;
	$("#jmlpoin-1-1").html(iki11);
	$("#jmlpoin-1-2").html(iki12);
	$("#jmlpoin-2-1").html(iki21);
	$("#jmlpoin-2-2").html(iki22);
	$("#jmlxpoin1").html(tot1.toFixed(4));
	$("#jmlxpoin2").html(tot2.toFixed(2));
	$("#totalpoin").html("Keterangan: Total Poin untuk perhitungan IKI adalah: " + totfinal.toFixed(4));
}

$("#navbar").on("click", ".bcptgl", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tgl = $(this).html();
	console.log(tgl)
	$.post("getbcpmonth", {
		token: token,
		tgl: tgl,
		email: $("#email").val()
	}, function(data){
		var js = JSON.parse(data);
		// console.log("Jumlah interna adalah : " + JSON.stringify(js.data.data1))
		pieChart(js.data, tgl)
		$("#inputnocm").hide();
		$(".diagram").show();
		$("#tabeliki").html(js.modal).show();
		getSum()
		$("div.tabtitle").html("Tabel IKI " + tgl);
		$("#tabelutama").html(js.script);
	})
});

$("#navbar").on("click", ".createpdf", function(e){
	e.preventDefault();
	var token = localStorage.getItem("token");
	var tgl = $(this).html();
	var nama = "dr. " + $("#dokter").html();
	var email = $("#email").val();
	$("#emailpdf").val(email);
	$("#namapdf").val(nama);
	$("#tglpdf").val(tgl);
	$("#tokenpdf").val(token);
	$("#getpdf").submit();

});

var pieChart = function(list, tgl){


	google.charts.load('current', {packages: ['corechart', 'bar']});
	google.charts.setOnLoadCallback(drawChart);
	google.charts.setOnLoadCallback(barChart);
	function drawChart(){
		var data = google.visualization.arrayToDataTable([
          ['Bagian', 'Jumlah Pasien'],
          ['Interna', list.data0.interna],
          ['Bedah', list.data0.bedah],
          ['Anak',  list.data0.anak],
          ['OBGYN', list.data0.obgyn],
          ['Saraf', list.data0.saraf],
          ['Anestesi', list.data0.anes],
          ['Psikiatri', list.data0.psik],
          ['THT', list.data0.tht],
          ['Kulit dan Kelamin', list.data0.kulit],
          ['Jantung', list.data0.jant],
          ['Umum', list.data0.um],
          ['Mata', list.data0.mata],
          ['MOD', list.data0.mod],
        ]);

        var options = {
		  title: 'Jumlah Pasien PerBagian',
		  width: 800,
		  height: 400,
        };

        var chart = new google.visualization.PieChart(document.getElementById('piechart'));

        chart.draw(data, options);
	}

	function barChart(){
		var data = new google.visualization.DataTable();
	    data.addColumn('number', 'Tanggal Jaga');
		data.addColumn('number', 'Jumlah IKI 1');
		data.addColumn('number', 'Jumlah IKI 2');
		// var title = [['Tanggal Jaga', 'Jumlah IKI 1', 'Jumlah IKI 2']]
		// var js = JSON.parse(list.data1)
		var isi = []
		for (i=0;i<31;i++){
			// var isi = array[awal.tgl, awal.iki1, awal.iki2]
			var awal = list.data1[i]
			console.log("Awal adalah : " + awal.tgl)
			data.addRows([
				[awal.tgl, awal.iki1, awal.iki2]
			])
		}

		
		// var data = new google.visualization.arrayToDataTable(title);
        var options = {
			title: 'Jumlah Pasien Pertanggal Jaga',
			width: 800,
			hAxis : {
				title: "Tanggal Jaga",
			},
			vAxis: {
				title: "Jumlah Pasien"
			}
        };
		var chart = new google.visualization.ColumnChart(document.getElementById('barchart'));

        chart.draw(data, options);

	}
}


});


