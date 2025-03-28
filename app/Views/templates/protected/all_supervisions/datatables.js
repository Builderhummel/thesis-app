$(document).ready(function () {
    let now = new Date();
    let formattedDate = now.getFullYear() + "_" +
                        ("0" + (now.getMonth() + 1)).slice(-2) + "_" +
                        ("0" + now.getDate()).slice(-2) + "_" +
                        ("0" + now.getHours()).slice(-2) + "_" +
                        ("0" + now.getMinutes()).slice(-2);

    $('#all_supervisions').DataTable({
        dom: '<"d-flex justify-content-between align-items-center mb-2"lfB>rtip',
        lengthMenu: [[10, 25, 50, -1], [10, 25, 50, "All"]],
        buttons: [
            {
                extend: 'csv',
                text: '<i class="fas fa-file-csv"></i> CSV',
                className: 'btn btn-success btn-sm',
                filename: formattedDate + '_all_supervisions',
                exportOptions: {
                    columns: ':not(:last-child)'
                }
            },
            {
                extend: 'excel',
                text: '<i class="fas fa-file-excel"></i> Excel',
                className: 'btn btn-success btn-sm',
                filename: formattedDate + '_all_supervisions',
                exportOptions: {
                    columns: ':not(:last-child)'
                }
            },
            {
                extend: 'pdf',
                text: '<i class="fas fa-file-pdf"></i> PDF',
                className: 'btn btn-danger btn-sm',
                filename: formattedDate + '_all_supervisions',
                exportOptions: {
                    columns: ':not(:last-child)'
                }
            },
            {
                extend: 'print',
                text: '<i class="fas fa-print"></i> Print',
                className: 'btn btn-primary btn-sm',
                filename: formattedDate + '_all_supervisions',
                exportOptions: {
                    columns: ':not(:last-child)'
                }
            }
        ]
    });
});
