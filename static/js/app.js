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
            
            // Handle metadata display
            displayMetadata(data.metadata);
            
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
    
    // Function to display metadata
    function displayMetadata(metadata) {
        const noMetadataElement = document.getElementById('noMetadata');
        const metadataContentElement = document.getElementById('metadataContent');
        
        // If no metadata is available
        if (!metadata || Object.keys(metadata).length === 0) {
            noMetadataElement.classList.remove('d-none');
            metadataContentElement.classList.add('d-none');
            return;
        }
        
        // Hide the no metadata message and show the content
        noMetadataElement.classList.add('d-none');
        metadataContentElement.classList.remove('d-none');
        
        // Display confidence score
        document.getElementById('confidenceScore').textContent = 
            metadata.confidenceScore ? 
            (metadata.confidenceScore * 100).toFixed(1) + '%' : 
            'Not available';
        
        // Handle time periods
        const timePeriodsElement = document.getElementById('timePeriods');
        const timePeriodsListElement = document.getElementById('timePeriodslist');
        
        if (metadata.timePeriods && metadata.timePeriods.length > 0) {
            timePeriodsElement.classList.remove('d-none');
            timePeriodsListElement.innerHTML = '';
            
            metadata.timePeriods.forEach(period => {
                const yearRange = period.startYear < 0 ? 
                    Math.abs(period.startYear) + ' BCE to ' + 
                    (period.endYear < 0 ? Math.abs(period.endYear) + ' BCE' : period.endYear + ' CE') :
                    period.startYear + ' CE to ' + period.endYear + ' CE';
                    
                const li = document.createElement('li');
                li.innerHTML = `<strong>${period.name}</strong>: ${yearRange}<br><span class="text-muted small">${period.description}</span>`;
                timePeriodsListElement.appendChild(li);
            });
        } else {
            timePeriodsElement.classList.add('d-none');
        }
        
        // Handle regions
        const regionsElement = document.getElementById('regions');
        const regionsListElement = document.getElementById('regionsList');
        
        if (metadata.regions && metadata.regions.length > 0) {
            regionsElement.classList.remove('d-none');
            regionsListElement.innerHTML = '';
            
            metadata.regions.forEach(region => {
                const modernAreas = region.modernAreas ? 
                    `<span class="text-muted small">Modern-day: ${region.modernAreas.join(', ')}</span>` : '';
                    
                const li = document.createElement('li');
                li.innerHTML = `<strong>${region.name}</strong><br>${region.description}<br>${modernAreas}`;
                regionsListElement.appendChild(li);
            });
        } else {
            regionsElement.classList.add('d-none');
        }
        
        // Handle cultural context
        const culturalContextElement = document.getElementById('culturalContext');
        const culturalContextTextElement = document.getElementById('culturalContextText');
        
        if (metadata.culturalContext && metadata.culturalContext.length > 0) {
            culturalContextElement.classList.remove('d-none');
            culturalContextTextElement.textContent = metadata.culturalContext.join(', ');
        } else {
            culturalContextElement.classList.add('d-none');
        }
        
        // Handle material context
        const materialContextElement = document.getElementById('materialContext');
        const materialContextTextElement = document.getElementById('materialContextText');
        
        if (metadata.materialContext && metadata.materialContext.length > 0) {
            materialContextElement.classList.remove('d-none');
            materialContextTextElement.textContent = metadata.materialContext.join(', ');
        } else {
            materialContextElement.classList.add('d-none');
        }
        
        // Handle historical events
        const historicalEventsElement = document.getElementById('historicalEvents');
        const historicalEventsListElement = document.getElementById('historicalEventsList');
        
        if (metadata.historicalEvents && metadata.historicalEvents.length > 0) {
            historicalEventsElement.classList.remove('d-none');
            historicalEventsListElement.innerHTML = '';
            
            metadata.historicalEvents.forEach(event => {
                const yearInfo = event.year ? ` (${event.year < 0 ? Math.abs(event.year) + ' BCE' : event.year + ' CE'})` : '';
                const li = document.createElement('li');
                li.innerHTML = `<strong>${event.name}${yearInfo}</strong> - ${event.eventType}<br><span class="text-muted small">${event.description}</span>`;
                historicalEventsListElement.appendChild(li);
            });
        } else {
            historicalEventsElement.classList.add('d-none');
        }
    }
});
