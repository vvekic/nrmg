using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Text;
using System.Windows.Forms;
using System.IO;
using Heroes5RandomMap.DataProc;

namespace Heroes5TerrainEditor
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
        }

        private TerrainData aTerrainData = new TerrainData(72);


        private void DisplayOpenBinFileDialog(object sender, EventArgs e)
        {
            OpenFileDialog AFileDialog = new OpenFileDialog();
            AFileDialog.InitialDirectory =  Environment.GetFolderPath(Environment.SpecialFolder.Desktop);
            AFileDialog.Filter = "地形文件|*.bin|All|*.*";
            DialogResult result = AFileDialog.ShowDialog();
            if (result == DialogResult.OK)
            {
                OpenBinFile(AFileDialog.FileName);
            }
        }
        private void DisplaySaveBinFileDialog(object sender, EventArgs e)
        {
            SaveFileDialog AFileDialog = new SaveFileDialog();
            AFileDialog.InitialDirectory = Environment.GetFolderPath(Environment.SpecialFolder.Desktop);
            AFileDialog.Filter = "地形文件|*.bin|All|*.*";
            DialogResult result = AFileDialog.ShowDialog();
            if (result == DialogResult.OK)
            {
                SaveBinFile(AFileDialog.FileName);
            }
        }

        private int OpenBinFile(string stFileName)
        {
            aTerrainData = new TerrainData(stFileName);
            //清理显示区
            label3.Text = "Texture";

            //显示平面表和绘图
            ShowTerrainDataBlockInfo(((IExportGridInfo)aTerrainData.nPlateau));

            //修改Texture的菜单和下拉列表
            DataTable dtTextureList = new DataTable();
            dtTextureList.Columns.Add("Display");
            dtTextureList.Columns.Add("Value");

            textureToolStripMenuItem.DropDownItems.Clear();
            for (int i=0;i<aTerrainData.nLayer;i++)
            {
                ToolStripMenuItem tsmiTemp = new ToolStripMenuItem(i.ToString(), null, new System.EventHandler(this.menuShowTextureInfo));
                textureToolStripMenuItem.DropDownItems.Add(tsmiTemp);

                DataRow drTemp = dtTextureList.NewRow();
                drTemp["Display"] = aTerrainData.nTextureList[i].stLayerName;
                drTemp["Value"] = i;
                dtTextureList.Rows.Add(drTemp);
            }

            comboBoxLayerNumber.DataSource = dtTextureList;
            comboBoxLayerNumber.DisplayMember = "Display";
            comboBoxLayerNumber.ValueMember = "Value";

            return 1;
        }
        private void SaveBinFile(string stFileName)
        {
            aTerrainData.Save(stFileName);
            MessageBox.Show("Saved", "Prompt");
        }
        private void NewBinFile(UInt32 nBaseLength)
        {
            aTerrainData = new TerrainData(nBaseLength);
        }

        private void ShowTerrainDataBlockInfo(IExportGridInfo aIExport)
        {
            textBoxGrid.Lines = aIExport.ExportText();
            pictureBoxPlateau.Image = aIExport.ExportBMP(pictureBoxPlateau.Size.Width, pictureBoxPlateau.Size.Height);
        }
        private void menuShowTextureInfo(object sender, EventArgs e)
        {
            int nLayerSerialNumber = Convert.ToInt32(((ToolStripMenuItem)sender).Text);
            ShowTextureInfo(nLayerSerialNumber);
        }
        private void ShowTextureInfo(int nLayerSerialNumber)
        {
            TerrainLayerBlock mlbTemp = aTerrainData.nTextureList[nLayerSerialNumber];
            label3.Text = mlbTemp.stTexturePath;
            ShowTerrainDataBlockInfo((IExportGridInfo)mlbTemp);
        }

        private void heightToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nHeight);
        }
        private void plateauToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nPlateau);
        }
        private void unknow1ToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nRamp);
        }
        private void unknow2ToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nWater);
        }
        private void passibleToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nPassible);
        }
        private void unknow4ToolStripMenuItem_Click(object sender, EventArgs e)
        {
            ShowTerrainDataBlockInfo((IExportGridInfo)aTerrainData.nUnknown4);
        }

        private void comboBoxLayerNumber_SelectionChangeCommitted(object sender, EventArgs e)
        {
            ShowTextureInfo((int)((ComboBox)sender).SelectedIndex);
        }

        private void newToolStripMenuItem_Click(object sender, EventArgs e)
        {
            NewBinFile((UInt32)136);
            //DisplaySaveBinFileDialog(null, null);
        }
    }
}
