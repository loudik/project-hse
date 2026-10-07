import ApiService from './ApiService'

export async function apiCreateVesselApplication(data) {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications',
        method: 'post',
        data,
    })
}

export async function apiUpdateVesselApplication(id, data) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}`,
        method: 'patch',
        data,
    })
}

export async function apiSubmitVesselApplication(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/submit`,
        method: 'post',
    })
}

export async function apiListVesselApplications() {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications',
        method: 'get',
    })
}

export async function apiGetVesselApplication(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}`,
        method: 'get',
    })
}

export async function apiLookupVesselApplication(number) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/lookup?number=${encodeURIComponent(number)}`,
        method: 'get',
    })
}

export async function apiWithdrawVesselApplication(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/withdraw`,
        method: 'post',
    })
}

export async function apiReopenVesselApplication(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/reopen`,
        method: 'post',
    })
}

export async function apiUploadVesselDocument(applicationId, formData) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${applicationId}/documents`,
        method: 'post',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' },
    })
}

export async function apiListVesselDocuments(applicationId) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${applicationId}/documents`,
        method: 'get',
    })
}

export async function apiDownloadVesselDocument(applicationId, docId) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${applicationId}/documents/${docId}/download`,
        method: 'get',
    })
}

export async function apiListVesselApplicationsForReview(status) {
    return ApiService.fetchDataWithAxios({
        url: status ? `/vessel-applications/review?status=${encodeURIComponent(status)}` : '/vessel-applications/review',
        method: 'get',
    })
}

export async function apiGetVesselApplicationForReview(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/review/${id}`,
        method: 'get',
    })
}

export async function apiListVesselDocumentsForReview(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/review/${id}/documents`,
        method: 'get',
    })
}

export async function apiDecideVesselApplication(id, status) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/decision`,
        method: 'patch',
        data: { status },
    })
}

export async function apiGetVesselApplicationHistory(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/history`,
        method: 'get',
    })
}

export async function apiListHSEOfficers() {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications/hse-officers',
        method: 'get',
    })
}

export async function apiAssignHSEOfficer(id, hseOfficerId) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/assign`,
        method: 'patch',
        data: { hseOfficerId },
    })
}

export async function apiGetAuthorisationLetter(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/authorisation`,
        method: 'get',
    })
}

export async function apiIssueAuthorisationLetter(id, formData) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/authorisation`,
        method: 'post',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' },
    })
}

export async function apiDownloadAuthorisationLetter(id, format = 'docx') {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/authorisation/download`,
        method: 'get',
        params: { format },
    })
}

export async function apiDownloadVesselOutcomeDocument(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/outcome-document/download`,
        method: 'get',
    })
}

export async function apiViewVesselOutcomeDocument(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/outcome-document/view`,
        method: 'get',
    })
}

export async function apiRequestExtension(id, requestedEntryDateTo, reason) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/extension`,
        method: 'post',
        data: { requestedEntryDateTo, reason },
    })
}

export async function apiListExtensionRequestsForApplication(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/extension-requests`,
        method: 'get',
    })
}

export async function apiListExtensionRequests(status = 'Pending') {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications/extension-requests',
        method: 'get',
        params: { status },
    })
}

export async function apiDecideExtension(id, status, rejectionReason) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/extension-requests/${id}/decision`,
        method: 'patch',
        data: { status, rejectionReason },
    })
}

export async function apiSubstituteVessel(originalApplicationNumber) {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications/substitute',
        method: 'post',
        data: { originalApplicationNumber },
    })
}

export async function apiListStaffUsers() {
    return ApiService.fetchDataWithAxios({
        url: '/vessel-applications/staff-users',
        method: 'get',
    })
}

export async function apiRequestStaffReview(id, staffIds) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/review-requests`,
        method: 'post',
        data: { staffIds },
    })
}

export async function apiListReviewRequests(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/review-requests`,
        method: 'get',
    })
}

export async function apiAddReviewComment(id, comment) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/comments`,
        method: 'post',
        data: { comment },
    })
}

export async function apiListReviewComments(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/comments`,
        method: 'get',
    })
}

export async function apiNotifyANPReady(id) {
    return ApiService.fetchDataWithAxios({
        url: `/vessel-applications/${id}/notify-anp-ready`,
        method: 'post',
    })
}