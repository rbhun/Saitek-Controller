#include "StdAfx.h"
#include "FipDlg.h"
#include "RawImage.h"

#include <atlimage.h>

// define a filter used by the file open dialog
#define FILE_OPEN_FILTER _T("JPEG (*.jpg)\0*.jpg\0Bitmap (*.bmp)\0*.bmp\0All Files (*.*)\0*.*\0")
namespace
{
	static LPCTSTR s_ErrorToString(HRESULT hr)
	{
		switch (hr)
		{
		case S_OK:				return _T("S_OK");
		case E_FAIL:			return _T("E_FAIL");
		case E_HANDLE:			return _T("E_HANDLE");
		case E_INVALIDARG:		return _T("E_INVALIDARG");
		case E_BUFFERTOOSMALL:	return _T("E_BUFFERTOOSMALL");
		default:				return _T("Unknown");
		}
	}
	static void s_RenderImage(HDC hdc, LPCTSTR tsz)
	{
		CImage image;
		HRESULT hr = image.Load(tsz);
		if (SUCCEEDED(hr))
		{
			int old = SetStretchBltMode(hdc, COLORONCOLOR);
			image.StretchBlt(hdc, 0, 0, 320, 240, SRCCOPY);
			SetStretchBltMode(hdc, old);
		} 
	}
}
CFipDlg::CFipDlg(CDirectOutput& directoutput, void* hDevice) : m_directoutput(directoutput), m_device(hDevice), m_leftscroll(0), m_rightscroll(0), m_init(false)
{
}
CFipDlg::~CFipDlg()
{
}
LRESULT CFipDlg::OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/)
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
	hr = m_directoutput.AddPage(m_device, 1, L"FIP Test Page", FLAG_SET_AS_ACTIVE);
	if (FAILED(hr))
	{
		// flag this failure
		DisplayErrorMessage(_T("DirectOutput_AddPage failed with error "), hr);
		return -1;
	}

	m_init = true;
	// Trigger a redraw on the device
	//OnPageChanged(m_device, 1, true, this);

	return 0;
}
LRESULT CFipDlg::OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
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
LRESULT CFipDlg::OnBrowse1(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	// display a dialog to browse for a bmp/jpg file
	CFileDialog file(TRUE, 0, 0, 4|2, FILE_OPEN_FILTER);
	if (file.DoModal() == IDOK)
	{
		// update the edit box
		::SendMessage(GetDlgItem(IDC_EDIT1), WM_SETTEXT, 0, (LPARAM)file.m_szFileName);
	}
	return 0;
}
LRESULT CFipDlg::OnDisplay1(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	if (::IsWindow(GetDlgItem(IDC_EDIT1)) && m_init)
	{
		// get the edit box's data
		TCHAR tszBuffer[1024] = { 0 };
		::SendMessage(GetDlgItem(IDC_EDIT1), WM_GETTEXT, (WPARAM)sizeof(tszBuffer)/sizeof(tszBuffer[0]), (LPARAM)tszBuffer);

		HRESULT hr;

		// set this file on the device
		// Note: the file must be the correct size
		CT2W wszBuffer(tszBuffer);
		const DWORD length( wcslen(wszBuffer) );

		// NOTE: a zero length file name is not valid
		if (length > 0)
		{
			CRawImage img;
			HDC hdc = img.BeginPaint();
			s_RenderImage(hdc, tszBuffer);
			img.EndPaint();

			hr = m_directoutput.SetImage(m_device, 1, 0, 320*240*3, img.Buffer());
			if (FAILED(hr))
			{
				// flag this failure
				// Note: E_PAGENOTACTIVE is not flagged
				if (hr != E_PAGENOTACTIVE)
				{
					DisplayErrorMessage(_T("DirectOutput_SetImageFromFile failed with error "), hr);
				}
			}
		} else
		{
			// Blank the screen by displaying a null image
			hr = m_directoutput.SetImage(m_device, 1, 0, 0, NULL);
		}
		if (FAILED(hr))
		{
			// flag this failure
			// Note: E_PAGENOTACTIVE is not flagged
			if (hr != E_PAGENOTACTIVE)
			{
				DisplayErrorMessage(_T("DirectOutput_SetImageFromFile failed with error "), hr);
			}
		}
	}
	return 0;
}
LRESULT CFipDlg::OnBrowse2(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	// display a dialog to browse for a bmp/jpg file
	CFileDialog file(TRUE, 0, 0, 4|2, FILE_OPEN_FILTER);
	if (file.DoModal() == IDOK)
	{
		// update the edit box
		::SendMessage(GetDlgItem(IDC_EDIT2), WM_SETTEXT, 0, (LPARAM)file.m_szFileName);
	}
	return 0;
}
LRESULT CFipDlg::OnDisplay2(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	if (::IsWindow(GetDlgItem(IDC_EDIT2)) && m_init)
	{
		// get the edit box's data
		TCHAR tszBuffer[1024] = { 0 };
		::SendMessage(GetDlgItem(IDC_EDIT2), WM_GETTEXT, (WPARAM)sizeof(tszBuffer)/sizeof(tszBuffer[0]), (LPARAM)tszBuffer);

		HRESULT hr;

		// NOTE: an empty string is not a valid argument
		if (_tcslen(tszBuffer) > 0)
		{
			// set this file on the device
			// Note: the file must be the correct size
			CRawImage img;
			HDC hdc = img.BeginPaint();
			s_RenderImage(hdc, tszBuffer);
			img.EndPaint();

			hr = m_directoutput.SetImage(m_device, 1, 0, 320*240*3, img.Buffer());
		} else
		{
			// Blank the screen by displaying a null image
			hr = m_directoutput.SetImage(m_device, 1, 0, 0, NULL);
		}
		if (FAILED(hr))
		{
			// flag this failure
			// Note: E_PAGENOTACTIVE is not flagged
			if (hr != E_PAGENOTACTIVE)
			{
				DisplayErrorMessage(_T("DirectOutput_SetImageFromFile failed with error "), hr);
			}
		}
	}
	return 0;
}
/*static*/ void __stdcall CFipDlg::OnPageChanged(void* hDevice, DWORD dwPage, bool bSetActive, void* pCtxt)
{
	CFipDlg* pThis = (CFipDlg*)pCtxt;
	if (bSetActive)
	{
		// resend this page data
		BOOL bHandled(FALSE);
		pThis->OnDisplay1(BN_CLICKED, IDC_BUTTON2, pThis->GetDlgItem(IDC_BUTTON2), bHandled);
	}
}
/*static*/ void __stdcall CFipDlg::OnSoftButtonChanged(void* hDevice, DWORD dwButtons, void* pCtxt)
{
	// update the display on the window
	CFipDlg* pThis = (CFipDlg*)pCtxt;
	// update the scroll positions
	if (dwButtons & SoftButton_Left)
		++pThis->m_leftscroll;
	else if (dwButtons & SoftButton_Right)
		--pThis->m_leftscroll;
	if (dwButtons & SoftButton_Up)
		++pThis->m_rightscroll;
	else if (dwButtons & SoftButton_Down)
		--pThis->m_rightscroll;
	// update the display
	TCHAR tszBuffer[256] = { 0 };
	_sntprintf_s(tszBuffer, sizeof(tszBuffer)/sizeof(tszBuffer[0]), sizeof(tszBuffer)/sizeof(tszBuffer[0]), _T("Buttons = %08X (%d) (%d)\n"), dwButtons, pThis->m_leftscroll, pThis->m_rightscroll);
	::SendMessage(::GetDlgItem(pThis->m_hWnd, IDC_BUTTON_TEXT), WM_SETTEXT, 0, (LPARAM)tszBuffer);
}
void CFipDlg::DisplayErrorMessage(LPCTSTR tszMsg, HRESULT hr)
{
	TCHAR tszBuffer[1024] = { 0 };
	_sntprintf_s(tszBuffer, sizeof(tszBuffer)/sizeof(tszBuffer[0]), sizeof(tszBuffer)/sizeof(tszBuffer[0]), _T("%s %08X %s\n"), tszMsg, hr, s_ErrorToString(hr));
	MessageBox(tszBuffer, _T("Test.exe - FipDlg"), MB_ICONERROR);
}
