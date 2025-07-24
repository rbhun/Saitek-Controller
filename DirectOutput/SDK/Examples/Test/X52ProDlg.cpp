#include "StdAfx.h"
#include "X52ProDlg.h"

namespace
{
	static LPCTSTR s_ErrorToString(HRESULT hr)
	{
		switch (hr)
		{
		case S_OK:			return _T("S_OK");
		case E_FAIL:		return _T("E_FAIL");
		case E_HANDLE:		return _T("E_HANDLE");
		case E_INVALIDARG:	return _T("E_INVALIDARG");
		default:			return _T("Unknown");
		}
	}
	static DWORD s_GetStringIdFromControlID(WORD wID)
	{
		switch (wID)
		{
		case IDC_EDIT1:		return 0;
		case IDC_EDIT2:		return 1;
		case IDC_EDIT3:		return 2;
		default:			return 0;
		}
	}
	static DWORD s_GetLedIdFromControlID(WORD wID)
	{
		switch (wID)
		{
		case IDC_CHECK1:	return 0;
		case IDC_CHECK2:	return 1;
		case IDC_CHECK3:	return 2;
		case IDC_CHECK4:	return 3;
		case IDC_CHECK5:	return 4;
		case IDC_CHECK6:	return 5;
		case IDC_CHECK7:	return 6;
		case IDC_CHECK8:	return 7;
		case IDC_CHECK9:	return 8;
		case IDC_CHECK10:	return 9;
		case IDC_CHECK11:	return 10;
		case IDC_CHECK12:	return 11;
		case IDC_CHECK13:	return 12;
		case IDC_CHECK14:	return 13;
		case IDC_CHECK15:	return 14;
		case IDC_CHECK16:	return 15;
		case IDC_CHECK17:	return 16;
		case IDC_CHECK18:	return 17;
		case IDC_CHECK19:	return 18;
		case IDC_CHECK20:	return 19;
		default:			return 0;
		}
	}
}
CX52ProDlg::CX52ProDlg(CDirectOutput& directoutput, void* hDevice) : m_directoutput(directoutput), m_device(hDevice), m_scrollpos(0)
{
}
LRESULT CX52ProDlg::OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/)
{
	HRESULT hr;

	// Register a callback that gets called when the page changes
	hr = m_directoutput.RegisterPageCallback(m_device, OnPageChanged, this);
	if (FAILED(hr))
	{
		// flag this failure
		DisplayErrorMessage(_T("DirectOutput_RegisterPageCallback failed with error "), hr);
		return -1;
	}

	// Register a callback that gets called when the soft buttons get changed
	hr = m_directoutput.RegisterSoftButtonCallback(m_device, OnSoftButtonChanged, this);
	if (FAILED(hr))
	{
		// flag this failure
		DisplayErrorMessage(_T("DirectOutput_RegisterSoftButtonCallback failed with error "), hr);
		return -1;
	}

	// Add a page to the device (Page 1)
	// Flag the page to be activated when created
	// NOTE: This will NOT call the OnPageChanged callback
	hr = m_directoutput.AddPage(m_device, 1, L"X52Pro Test Page", FLAG_SET_AS_ACTIVE);
	if (FAILED(hr))
	{
		// flag this failure
		DisplayErrorMessage(_T("DirectOutput_AddPage failed with error "), hr);
		return -1;
	}

	// Trigger a redraw on the device
	OnPageChanged(m_device, 1, true, this);

	return 0;
}
LRESULT CX52ProDlg::OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	HRESULT hr;

	// Remove page 1
	hr = m_directoutput.RemovePage(m_device, 1);
	if (FAILED(hr))
	{
		// flag this failure
		DisplayErrorMessage(_T("DirectOutput_RemovePage failed with error "), hr);
	}

	EndDialog(wID);
	return 0;
}
LRESULT CX52ProDlg::OnEditChanged(WORD /*wNotifyCode*/, WORD wID, HWND hWndCtl, BOOL& /*bHandled*/)
{
	TCHAR tszBuffer[1024] = { 0 };
	::SendMessage(hWndCtl, WM_GETTEXT, (WPARAM)sizeof(tszBuffer)/sizeof(tszBuffer[0]), (LPARAM)tszBuffer);
	const DWORD id( s_GetStringIdFromControlID(wID) );

	HRESULT hr;

	CT2W wszBuffer(tszBuffer);
	hr = m_directoutput.SetString(m_device, 1, id, wcslen(wszBuffer), wszBuffer);
	if (FAILED(hr))
	{
		// flag this failure
		// Note: E_PAGENOTACTIVE is not flagged
		if (hr != E_PAGENOTACTIVE)
		{
			DisplayErrorMessage(_T("DirectOutput_SetString failed with error "), hr);
		}
	}

	return 0;
}
LRESULT CX52ProDlg::OnCheckChanged(WORD /*wNotifyCode*/, WORD wID, HWND hWndCtl, BOOL& /*bHandled*/)
{
	long lret = ::SendMessage(hWndCtl, BM_GETCHECK, 0, 0);
	const DWORD value( lret == BST_CHECKED ? 1 : 0 );
	const DWORD id( s_GetLedIdFromControlID(wID) );

	HRESULT hr;

	hr = m_directoutput.SetLed(m_device, 1, id, value);
	if (FAILED(hr))
	{
		// flag this failure
		// Note: E_PAGENOTACTIVE is not flagged
		if (hr != E_PAGENOTACTIVE)
		{
			DisplayErrorMessage(_T("DirectOutput_SetLed failed with error "), hr);
		}
	}
	return 0;
}
/*static*/ void __stdcall CX52ProDlg::OnPageChanged(void* hDevice, DWORD dwPage, bool bSetActive, void* pCtxt)
{
	CX52ProDlg* pThis = (CX52ProDlg*)pCtxt;
	if (bSetActive)
	{
		// resend this page data
		BOOL bHandled(FALSE);
		pThis->OnEditChanged(EN_CHANGE, IDC_EDIT1, pThis->GetDlgItem(IDC_EDIT1), bHandled);
		pThis->OnEditChanged(EN_CHANGE, IDC_EDIT2, pThis->GetDlgItem(IDC_EDIT2), bHandled);
		pThis->OnEditChanged(EN_CHANGE, IDC_EDIT3, pThis->GetDlgItem(IDC_EDIT3), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK1, pThis->GetDlgItem(IDC_CHECK1), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK2, pThis->GetDlgItem(IDC_CHECK2), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK3, pThis->GetDlgItem(IDC_CHECK3), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK4, pThis->GetDlgItem(IDC_CHECK4), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK5, pThis->GetDlgItem(IDC_CHECK5), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK6, pThis->GetDlgItem(IDC_CHECK6), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK7, pThis->GetDlgItem(IDC_CHECK7), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK8, pThis->GetDlgItem(IDC_CHECK8), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK9, pThis->GetDlgItem(IDC_CHECK9), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK10, pThis->GetDlgItem(IDC_CHECK10), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK11, pThis->GetDlgItem(IDC_CHECK11), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK12, pThis->GetDlgItem(IDC_CHECK12), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK13, pThis->GetDlgItem(IDC_CHECK13), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK14, pThis->GetDlgItem(IDC_CHECK14), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK15, pThis->GetDlgItem(IDC_CHECK15), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK16, pThis->GetDlgItem(IDC_CHECK16), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK17, pThis->GetDlgItem(IDC_CHECK17), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK18, pThis->GetDlgItem(IDC_CHECK18), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK19, pThis->GetDlgItem(IDC_CHECK19), bHandled);
		pThis->OnCheckChanged(BN_CLICKED, IDC_CHECK20, pThis->GetDlgItem(IDC_CHECK20), bHandled);
	}
}
/*static*/ void __stdcall CX52ProDlg::OnSoftButtonChanged(void* hDevice, DWORD dwButtons, void* pCtxt)
{
	// update the display on the window
	CX52ProDlg* pThis = (CX52ProDlg*)pCtxt;
	// update the current position
	if (dwButtons & 0x00000002)
		++pThis->m_scrollpos;
	else if (dwButtons & 0x0000004)
		--pThis->m_scrollpos;
	// draw the string
	TCHAR tszBuffer[256] = { 0 };
	_sntprintf_s(tszBuffer, sizeof(tszBuffer)/sizeof(tszBuffer[0]), sizeof(tszBuffer)/sizeof(tszBuffer[0]), _T("Buttons = %08X (%d)\n"), dwButtons, pThis->m_scrollpos);
	::SendMessage(::GetDlgItem(pThis->m_hWnd, IDC_BUTTON_TEXT), WM_SETTEXT, 0, (LPARAM)tszBuffer);
}
void CX52ProDlg::DisplayErrorMessage(LPCTSTR tszMsg, HRESULT hr)
{
	TCHAR tszBuffer[1024] = { 0 };
	_sntprintf_s(tszBuffer, sizeof(tszBuffer)/sizeof(tszBuffer[0]), sizeof(tszBuffer)/sizeof(tszBuffer[0]), _T("%s %08X %s\n"), tszMsg, hr, s_ErrorToString(hr));
	MessageBox(tszBuffer, _T("Test.exe - X52ProDlg"), MB_ICONERROR);
}

