document.addEventListener('DOMContentLoaded', function() {
    // Get form elements
    const uploadForm = document.getElementById('uploadForm');
    const summarizeForm = document.getElementById('summarizeForm');
    
    // Get result display elements
    const loadingTranslation = document.getElementById('loadingTranslation');
    const translationResult = document.getElementById('translationResult');
    const noTranslation = document.getElementById('noTranslation');
    
    const loadingSummary = document.getElementById('loadingSummary');
    const summaryResult = document.getElementById('summaryResult');
    const noSummary = document.getElementById('noSummary');
    
    // Add event listener for the translation form
    uploadForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Show loading indicator
        loadingTranslation.classList.remove('d-none');
        translationResult.classList.add('d-none');
        noTranslation.classList.add('d-none');
        
        // Create FormData object
        const formData = new FormData(uploadForm);
        
        // Send AJAX request
        fetch('/api/translate', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(text || 'Failed to translate manuscript');
                });
            }
            return response.json();
        })
        .then(data => {
            // Hide loading indicator
            loadingTranslation.classList.add('d-none');
            
            // Populate result fields
            document.getElementById('originalScript').textContent = data.originalScript;
            document.getElementById('translatedText').textContent = data.translatedText;
            document.getElementById('summary').textContent = data.summary;
            document.getElementById('processedAt').textContent = data.processedAt;
            
            // Show result
            translationResult.classList.remove('d-none');
        })
        .catch(error => {
            // Hide loading indicator
            loadingTranslation.classList.add('d-none');
            
            // Show error
            alert('Error: ' + error.message);
            
            // Show no translation message
            noTranslation.classList.remove('d-none');
        });
    });
    
    // Add event listener for the summarization form
    summarizeForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Show loading indicator
        loadingSummary.classList.remove('d-none');
        summaryResult.classList.add('d-none');
        noSummary.classList.add('d-none');
        
        // Get text to summarize
        const text = document.getElementById('textToSummarize').value;
        
        // Send AJAX request
        fetch('/api/summarize', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ text: text })
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(text || 'Failed to summarize text');
                });
            }
            return response.json();
        })
        .then(data => {
            // Hide loading indicator
            loadingSummary.classList.add('d-none');
            
            // Populate result fields
            document.getElementById('directSummary').textContent = data.summary;
            document.getElementById('textLength').textContent = data.textLength;
            document.getElementById('summaryProcessedAt').textContent = data.processedAt;
            
            // Show result
            summaryResult.classList.remove('d-none');
        })
        .catch(error => {
            // Hide loading indicator
            loadingSummary.classList.add('d-none');
            
            // Show error
            alert('Error: ' + error.message);
            
            // Show no summary message
            noSummary.classList.remove('d-none');
        });
    });
    
    // Initialize feather icons if available
    if (typeof feather !== 'undefined') {
        feather.replace();
    }
});
