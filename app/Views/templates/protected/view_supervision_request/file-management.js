// Initialize event listeners when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    const fileInput = document.getElementById('fileInput');
    if (fileInput) {
        fileInput.addEventListener('change', uploadFile);
    }

    // Event delegation for delete buttons
    document.addEventListener('click', function(e) {
        if (e.target.classList.contains('btn-delete-file')) {
            const fuid = e.target.getAttribute('data-fuid');
            if (fuid) {
                deleteFile(fuid);
            }
        }
        
        // Handle category download buttons
        if (e.target.classList.contains('btn-download-category')) {
            const category = e.target.getAttribute('data-category');
            if (category) {
                downloadLatestByCategory(category);
            }
        }
    });
});

function uploadFile() {
    const form = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const categorySelect = document.getElementById('categorySelect');
    
    if (!fileInput.files || fileInput.files.length === 0) {
        return;
    }

    const formData = new FormData(form);
    
    fetch('/files/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('Error uploading file: ' + data.error);
        } else {
            location.reload(); // Reload page to show new file
        }
    })
    .catch(error => {
        alert('Error uploading file: ' + error);
    });
}

function deleteFile(fuid) {
    if (!confirm('Are you sure you want to delete this file?')) {
        return;
    }

    fetch('/files/delete?fuid=' + fuid, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert('Error deleting file: ' + data.error);
        } else {
            location.reload(); // Reload page to update file list
        }
    })
    .catch(error => {
        alert('Error deleting file: ' + error);
    });
}

function downloadLatestByCategory(category) {
    const btnGroup = document.querySelector('.btn-group[data-tuid]');
    if (!btnGroup) {
        alert('Unable to determine thesis ID');
        return;
    }

    const tuid = btnGroup.getAttribute('data-tuid');
    if (!tuid) {
        alert('Thesis ID not found');
        return;
    }

    // Download using category parameter
    window.location.href = `/files/download?tuid=${tuid}&category=${category}`;
}

function getCategoryName(category) {
    const names = {
        'transcript-of-records': 'Transcript of Records',
        'cv': 'CV',
        'final-thesis': 'Final Thesis'
    };
    return names[category] || category;
}
