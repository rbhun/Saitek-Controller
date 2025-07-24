#include "stdafx.h"
#include "resource.h"

#include "MainDlg.h"
#include "FipDlg.h"
#include "X52ProDlg.h"

LRESULT CMainDlg::OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/)
{
	// center the dialog on the screen
	CenterWindow();

	// set icons
	HICON hIcon = (HICON)::LoadImage(_Module.GetResourceInstance(), MAKEINTRESOURCE(IDR_MAINFRAME), 
		IMAGE_ICON, ::GetSystemMetrics(SM_CXICON), ::GetSystemMetrics(SM_CYICON), LR_DEFAULTCOLOR);
	SetIcon(hIcon, TRUE);
	HICON hIconSmall = (HICON)::LoadImage(_Module.GetResourceInstance(), MAKEINTRESOURCE(IDR_MAINFRAME), 
		IMAGE_ICON, ::GetSystemMetrics(SM_CXSMICON), ::GetSystemMetrics(SM_CYSMICON), LR_DEFAULTCOLOR);
	SetIcon(hIconSmall, FALSE);

	// Initialize DirectOutput
	InitializeDirectOutput();

	return TRUE;
}

LRESULT CMainDlg::OnAppAbout(WORD /*wNotifyCode*/, WORD /*wID*/, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	CSimpleDialog<IDD_ABOUTBOX, FALSE> dlg;
	dlg.DoModal();
	return 0;
}

LRESULT CMainDlg::OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	EndDialog(wID);
	return 0;
}
LRESULT CMainDlg::OnListBoxDoubleClick(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/)
{
	long ret = ::SendMessage(GetDlgItem(IDC_LIST1), LB_GETCURSEL, 0, 0);
	void* hDevice = (void*)::SendMessage(GetDlgItem(IDC_LIST1), LB_GETITEMDATA, ret, 0);
	for (DeviceList::iterator it = m_devices.begin(); it != m_devices.end(); ++it)
	{
		if (*it == hDevice)
		{
			// what type of device is this?
			HRESULT hr;

			GUID typeguid;
			hr = m_directoutput.GetDeviceType(hDevice, &typeguid);
			if (FAILED(hr))
			{
				// flag this error
			}

			if (typeguid == DeviceType_X52Pro)
			{
				CX52ProDlg dlg(m_directoutput, hDevice);
				dlg.DoModal();
			}
			else if (typeguid == DeviceType_Fip)
			{
				CFipDlg dlg(m_directoutput, hDevice);
				dlg.DoModal();
			}

			break;
		}
	}
	return 0;
}

/*virtual*/ void CMainDlg::OnFinalMessage(HWND hWnd)
{
	// Cleanup DirectOutput
	m_directoutput.Deinitialize();
}

void CMainDlg::InitializeDirectOutput()
{
	HRESULT hr;
	
	// Initialize DirectOutput
	hr = m_directoutput.Initialize(L"Test");
	if (FAILED(hr))
	{
		// flag this error
	}

	// Register a callback to be called when a device is added or removed
	hr = m_directoutput.RegisterDeviceCallback(OnDeviceChanged, this);
	if (FAILED(hr))
	{
		// flag this error
	}

	// Enumerate all currently attached devices. This will call the device change callback
	hr = m_directoutput.Enumerate(OnEnumerateDevice, this);
	if (FAILED(hr))
	{
		// flag this error
	}

	UpdateListBox();
}
/*static*/ void __stdcall CMainDlg::OnEnumerateDevice(void* hDevice, void* pCtxt)
{
	CMainDlg* pThis = (CMainDlg*)pCtxt;
	pThis->m_devices.push_back(hDevice);
}
/*static*/ void __stdcall CMainDlg::OnDeviceChanged(void* hDevice, bool bAdded, void* pCtxt)
{
	CMainDlg* pThis = (CMainDlg*)pCtxt;
	if (bAdded)
	{
		// device has been added, add to list of devices
		{
			TCHAR tsz[1024];
			_sntprintf_s(tsz, sizeof(tsz)/sizeof(tsz[0]), sizeof(tsz)/sizeof(tsz[0]), _T("DeviceAdded(%p)\n"), hDevice);
			OutputDebugString(tsz);
		}
		pThis->m_devices.push_back(hDevice);
	}
	else
	{
		// device has been removed, remove from list of devices
		{
			TCHAR tsz[1024];
			_sntprintf_s(tsz, sizeof(tsz)/sizeof(tsz[0]), sizeof(tsz)/sizeof(tsz[0]), _T("DeviceRemoved(%p)\n"), hDevice);
			OutputDebugString(tsz);
		}
		for (DeviceList::iterator it = pThis->m_devices.begin(); it != pThis->m_devices.end(); ++it)
		{
			if (*it == hDevice)
			{
				pThis->m_devices.erase(it);
				break;
			}
		}
	}

	// update the list box control
	pThis->UpdateListBox();
}
void CMainDlg::UpdateListBox()
{
	// clear the listbox
	while (::SendMessage(GetDlgItem(IDC_LIST1), LB_DELETESTRING, 0, 0) > 0);
	// for each device, add an entry into the list box
	for (DeviceList::iterator it = m_devices.begin(); it != m_devices.end(); ++it)
	{
		GUID typeguid = { 0 };
		HRESULT hr;
		
		hr = m_directoutput.GetDeviceType(*it, &typeguid);
		if (FAILED(hr))
		{
			// flag this error
		}

		if (typeguid == DeviceType_X52Pro)
		{
			long ret = ::SendMessage(GetDlgItem(IDC_LIST1), LB_ADDSTRING, 0, (LPARAM)_T("X52Pro Device"));
			::SendMessage(GetDlgItem(IDC_LIST1), LB_SETITEMDATA, ret, (LPARAM)*it);
		}
		else if (typeguid == DeviceType_Fip)
		{
			long ret = ::SendMessage(GetDlgItem(IDC_LIST1), LB_ADDSTRING, 0, (LPARAM)_T("FIP Device"));
			::SendMessage(GetDlgItem(IDC_LIST1), LB_SETITEMDATA, ret, (LPARAM)*it);
		}
		else
		{
			long ret = ::SendMessage(GetDlgItem(IDC_LIST1), LB_ADDSTRING, 0, (LPARAM)_T("Unknown Device"));
			::SendMessage(GetDlgItem(IDC_LIST1), LB_SETITEMDATA, ret, (LPARAM)*it);
		}
	}
}
