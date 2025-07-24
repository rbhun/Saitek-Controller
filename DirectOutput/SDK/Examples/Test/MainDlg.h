#pragma once
#include "DirectOutputImpl.h"
#include <vector>

///
/// CMainDlg
///
/// Handles the Main Dialog, which enumerates attached devices
///
class CMainDlg : public CDialogImpl<CMainDlg>
{
	typedef std::vector<void*> DeviceList;

	///
	/// DirectOutput.dll Interface and collection of device handles
	///
	CDirectOutput	m_directoutput;
	DeviceList		m_devices;

	///
	/// Initialize the DirectOutput.dll communication
	///
	void InitializeDirectOutput();

	///
	/// Callbacks
	///
	static void __stdcall OnEnumerateDevice(void* hDevice, void* pCtxt);
	static void __stdcall OnDeviceChanged(void* hDevice, bool bAdded, void* pCtxt);

	///
	/// Helper Method that populates the listbox
	///
	void UpdateListBox();

public:
	///
	/// WTL Dialog Methods
	///
	enum { IDD = IDD_MAINDLG };

	BEGIN_MSG_MAP(CMainDlg)
		MESSAGE_HANDLER(WM_INITDIALOG, OnInitDialog)
		COMMAND_ID_HANDLER(ID_APP_ABOUT, OnAppAbout)
		COMMAND_ID_HANDLER(IDOK, OnClose)
		COMMAND_ID_HANDLER(IDCANCEL, OnClose)
		COMMAND_HANDLER(IDC_LIST1, LBN_DBLCLK, OnListBoxDoubleClick)
	END_MSG_MAP()

	LRESULT OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/);
	LRESULT OnAppAbout(WORD /*wNotifyCode*/, WORD /*wID*/, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnListBoxDoubleClick(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);

	virtual void OnFinalMessage(HWND hWnd);
};
