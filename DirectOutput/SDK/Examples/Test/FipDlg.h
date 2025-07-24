#pragma once
#include "resource.h"
#include "DirectOutputImpl.h"

///
/// CFipDlg
///
/// Handles the FIP Device Test Dialog
///
class CFipDlg : public CDialogImpl<CFipDlg>
{
public:
	CFipDlg(CDirectOutput& directoutput, void* hDevice);
	~CFipDlg();

	///
	/// Callbacks
	///
	static void __stdcall OnPageChanged(void* hDevice, DWORD dwPage, bool bSetActive, void* pCtxt);
	static void __stdcall OnSoftButtonChanged(void* hDevice, DWORD dwButtons, void* pCtxt);

private:
	///
	/// The Device this dialog controls
	///
	CDirectOutput&	m_directoutput;
	void*			m_device;
	int				m_leftscroll;
	int				m_rightscroll;
	bool			m_init;

	///
	/// Helper method to display an error message
	///
	void DisplayErrorMessage(LPCTSTR tszMsg, HRESULT hr);

public:
	/// 
	/// WTL Dialog Methods
	///
	enum { IDD = IDD_FIP };

	BEGIN_MSG_MAP(CFipDlg)
		MESSAGE_HANDLER(WM_INITDIALOG, OnInitDialog)
		COMMAND_ID_HANDLER(IDOK, OnClose)
		COMMAND_ID_HANDLER(IDCANCEL, OnClose)
		COMMAND_ID_HANDLER(IDC_BUTTON1, OnBrowse1)
		COMMAND_ID_HANDLER(IDC_BUTTON2, OnDisplay1)
		COMMAND_ID_HANDLER(IDC_BUTTON3, OnBrowse2)
		COMMAND_ID_HANDLER(IDC_BUTTON4, OnDisplay2)
	END_MSG_MAP()

	LRESULT OnInitDialog(UINT /*uMsg*/, WPARAM /*wParam*/, LPARAM /*lParam*/, BOOL& /*bHandled*/);
	LRESULT OnClose(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnBrowse1(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnDisplay1(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnBrowse2(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
	LRESULT OnDisplay2(WORD /*wNotifyCode*/, WORD wID, HWND /*hWndCtl*/, BOOL& /*bHandled*/);
};
