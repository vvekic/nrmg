namespace Heroes5TerrainEditor
{
    partial class Form1
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Form1));
            this.button1 = new System.Windows.Forms.Button();
            this.pictureBoxPlateau = new System.Windows.Forms.PictureBox();
            this.menuStripMain = new System.Windows.Forms.MenuStrip();
            this.fileToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.newToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.openToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.saveToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.viewToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.textureToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripMenuItem2 = new System.Windows.Forms.ToolStripMenuItem();
            this.heightToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.plateauToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.unknow1ToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.unknow2ToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.passibleToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.unknow4ToolStripMenuItem = new System.Windows.Forms.ToolStripMenuItem();
            this.textBoxGrid = new System.Windows.Forms.TextBox();
            this.label3 = new System.Windows.Forms.Label();
            this.comboBoxLayerNumber = new System.Windows.Forms.ComboBox();
            this.label1 = new System.Windows.Forms.Label();
            ((System.ComponentModel.ISupportInitialize)(this.pictureBoxPlateau)).BeginInit();
            this.menuStripMain.SuspendLayout();
            this.SuspendLayout();
            // 
            // button1
            // 
            this.button1.AccessibleDescription = null;
            this.button1.AccessibleName = null;
            resources.ApplyResources(this.button1, "button1");
            this.button1.BackgroundImage = null;
            this.button1.Font = null;
            this.button1.Name = "button1";
            this.button1.UseVisualStyleBackColor = true;
            this.button1.Click += new System.EventHandler(this.DisplayOpenBinFileDialog);
            // 
            // pictureBoxPlateau
            // 
            this.pictureBoxPlateau.AccessibleDescription = null;
            this.pictureBoxPlateau.AccessibleName = null;
            resources.ApplyResources(this.pictureBoxPlateau, "pictureBoxPlateau");
            this.pictureBoxPlateau.BackgroundImage = null;
            this.pictureBoxPlateau.Font = null;
            this.pictureBoxPlateau.ImageLocation = null;
            this.pictureBoxPlateau.Name = "pictureBoxPlateau";
            this.pictureBoxPlateau.TabStop = false;
            // 
            // menuStripMain
            // 
            this.menuStripMain.AccessibleDescription = null;
            this.menuStripMain.AccessibleName = null;
            resources.ApplyResources(this.menuStripMain, "menuStripMain");
            this.menuStripMain.BackgroundImage = null;
            this.menuStripMain.Font = null;
            this.menuStripMain.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.fileToolStripMenuItem,
            this.viewToolStripMenuItem});
            this.menuStripMain.Name = "menuStripMain";
            // 
            // fileToolStripMenuItem
            // 
            this.fileToolStripMenuItem.AccessibleDescription = null;
            this.fileToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.fileToolStripMenuItem, "fileToolStripMenuItem");
            this.fileToolStripMenuItem.BackgroundImage = null;
            this.fileToolStripMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.newToolStripMenuItem,
            this.openToolStripMenuItem,
            this.saveToolStripMenuItem});
            this.fileToolStripMenuItem.Name = "fileToolStripMenuItem";
            this.fileToolStripMenuItem.ShortcutKeyDisplayString = null;
            // 
            // newToolStripMenuItem
            // 
            this.newToolStripMenuItem.AccessibleDescription = null;
            this.newToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.newToolStripMenuItem, "newToolStripMenuItem");
            this.newToolStripMenuItem.BackgroundImage = null;
            this.newToolStripMenuItem.Name = "newToolStripMenuItem";
            this.newToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.newToolStripMenuItem.Click += new System.EventHandler(this.newToolStripMenuItem_Click);
            // 
            // openToolStripMenuItem
            // 
            this.openToolStripMenuItem.AccessibleDescription = null;
            this.openToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.openToolStripMenuItem, "openToolStripMenuItem");
            this.openToolStripMenuItem.BackgroundImage = null;
            this.openToolStripMenuItem.Name = "openToolStripMenuItem";
            this.openToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.openToolStripMenuItem.Click += new System.EventHandler(this.DisplayOpenBinFileDialog);
            // 
            // saveToolStripMenuItem
            // 
            this.saveToolStripMenuItem.AccessibleDescription = null;
            this.saveToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.saveToolStripMenuItem, "saveToolStripMenuItem");
            this.saveToolStripMenuItem.BackgroundImage = null;
            this.saveToolStripMenuItem.Name = "saveToolStripMenuItem";
            this.saveToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.saveToolStripMenuItem.Click += new System.EventHandler(this.DisplaySaveBinFileDialog);
            // 
            // viewToolStripMenuItem
            // 
            this.viewToolStripMenuItem.AccessibleDescription = null;
            this.viewToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.viewToolStripMenuItem, "viewToolStripMenuItem");
            this.viewToolStripMenuItem.BackgroundImage = null;
            this.viewToolStripMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.textureToolStripMenuItem,
            this.heightToolStripMenuItem,
            this.plateauToolStripMenuItem,
            this.unknow1ToolStripMenuItem,
            this.unknow2ToolStripMenuItem,
            this.passibleToolStripMenuItem,
            this.unknow4ToolStripMenuItem});
            this.viewToolStripMenuItem.Name = "viewToolStripMenuItem";
            this.viewToolStripMenuItem.ShortcutKeyDisplayString = null;
            // 
            // textureToolStripMenuItem
            // 
            this.textureToolStripMenuItem.AccessibleDescription = null;
            this.textureToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.textureToolStripMenuItem, "textureToolStripMenuItem");
            this.textureToolStripMenuItem.BackgroundImage = null;
            this.textureToolStripMenuItem.DropDownItems.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.toolStripMenuItem2});
            this.textureToolStripMenuItem.Name = "textureToolStripMenuItem";
            this.textureToolStripMenuItem.ShortcutKeyDisplayString = null;
            // 
            // toolStripMenuItem2
            // 
            this.toolStripMenuItem2.AccessibleDescription = null;
            this.toolStripMenuItem2.AccessibleName = null;
            resources.ApplyResources(this.toolStripMenuItem2, "toolStripMenuItem2");
            this.toolStripMenuItem2.BackgroundImage = null;
            this.toolStripMenuItem2.Name = "toolStripMenuItem2";
            this.toolStripMenuItem2.ShortcutKeyDisplayString = null;
            // 
            // heightToolStripMenuItem
            // 
            this.heightToolStripMenuItem.AccessibleDescription = null;
            this.heightToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.heightToolStripMenuItem, "heightToolStripMenuItem");
            this.heightToolStripMenuItem.BackgroundImage = null;
            this.heightToolStripMenuItem.Name = "heightToolStripMenuItem";
            this.heightToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.heightToolStripMenuItem.Click += new System.EventHandler(this.heightToolStripMenuItem_Click);
            // 
            // plateauToolStripMenuItem
            // 
            this.plateauToolStripMenuItem.AccessibleDescription = null;
            this.plateauToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.plateauToolStripMenuItem, "plateauToolStripMenuItem");
            this.plateauToolStripMenuItem.BackgroundImage = null;
            this.plateauToolStripMenuItem.Name = "plateauToolStripMenuItem";
            this.plateauToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.plateauToolStripMenuItem.Click += new System.EventHandler(this.plateauToolStripMenuItem_Click);
            // 
            // unknow1ToolStripMenuItem
            // 
            this.unknow1ToolStripMenuItem.AccessibleDescription = null;
            this.unknow1ToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.unknow1ToolStripMenuItem, "unknow1ToolStripMenuItem");
            this.unknow1ToolStripMenuItem.BackgroundImage = null;
            this.unknow1ToolStripMenuItem.Name = "unknow1ToolStripMenuItem";
            this.unknow1ToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.unknow1ToolStripMenuItem.Click += new System.EventHandler(this.unknow1ToolStripMenuItem_Click);
            // 
            // unknow2ToolStripMenuItem
            // 
            this.unknow2ToolStripMenuItem.AccessibleDescription = null;
            this.unknow2ToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.unknow2ToolStripMenuItem, "unknow2ToolStripMenuItem");
            this.unknow2ToolStripMenuItem.BackgroundImage = null;
            this.unknow2ToolStripMenuItem.Name = "unknow2ToolStripMenuItem";
            this.unknow2ToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.unknow2ToolStripMenuItem.Click += new System.EventHandler(this.unknow2ToolStripMenuItem_Click);
            // 
            // passibleToolStripMenuItem
            // 
            this.passibleToolStripMenuItem.AccessibleDescription = null;
            this.passibleToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.passibleToolStripMenuItem, "passibleToolStripMenuItem");
            this.passibleToolStripMenuItem.BackgroundImage = null;
            this.passibleToolStripMenuItem.Name = "passibleToolStripMenuItem";
            this.passibleToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.passibleToolStripMenuItem.Click += new System.EventHandler(this.passibleToolStripMenuItem_Click);
            // 
            // unknow4ToolStripMenuItem
            // 
            this.unknow4ToolStripMenuItem.AccessibleDescription = null;
            this.unknow4ToolStripMenuItem.AccessibleName = null;
            resources.ApplyResources(this.unknow4ToolStripMenuItem, "unknow4ToolStripMenuItem");
            this.unknow4ToolStripMenuItem.BackgroundImage = null;
            this.unknow4ToolStripMenuItem.Name = "unknow4ToolStripMenuItem";
            this.unknow4ToolStripMenuItem.ShortcutKeyDisplayString = null;
            this.unknow4ToolStripMenuItem.Click += new System.EventHandler(this.unknow4ToolStripMenuItem_Click);
            // 
            // textBoxGrid
            // 
            this.textBoxGrid.AccessibleDescription = null;
            this.textBoxGrid.AccessibleName = null;
            resources.ApplyResources(this.textBoxGrid, "textBoxGrid");
            this.textBoxGrid.BackgroundImage = null;
            this.textBoxGrid.Font = null;
            this.textBoxGrid.Name = "textBoxGrid";
            // 
            // label3
            // 
            this.label3.AccessibleDescription = null;
            this.label3.AccessibleName = null;
            resources.ApplyResources(this.label3, "label3");
            this.label3.Font = null;
            this.label3.Name = "label3";
            // 
            // comboBoxLayerNumber
            // 
            this.comboBoxLayerNumber.AccessibleDescription = null;
            this.comboBoxLayerNumber.AccessibleName = null;
            resources.ApplyResources(this.comboBoxLayerNumber, "comboBoxLayerNumber");
            this.comboBoxLayerNumber.BackgroundImage = null;
            this.comboBoxLayerNumber.Font = null;
            this.comboBoxLayerNumber.FormattingEnabled = true;
            this.comboBoxLayerNumber.Name = "comboBoxLayerNumber";
            this.comboBoxLayerNumber.SelectionChangeCommitted += new System.EventHandler(this.comboBoxLayerNumber_SelectionChangeCommitted);
            // 
            // label1
            // 
            this.label1.AccessibleDescription = null;
            this.label1.AccessibleName = null;
            resources.ApplyResources(this.label1, "label1");
            this.label1.Font = null;
            this.label1.Name = "label1";
            // 
            // Form1
            // 
            this.AccessibleDescription = null;
            this.AccessibleName = null;
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.None;
            resources.ApplyResources(this, "$this");
            this.BackgroundImage = null;
            this.Controls.Add(this.label1);
            this.Controls.Add(this.comboBoxLayerNumber);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.textBoxGrid);
            this.Controls.Add(this.pictureBoxPlateau);
            this.Controls.Add(this.button1);
            this.Controls.Add(this.menuStripMain);
            this.Icon = null;
            this.MainMenuStrip = this.menuStripMain;
            this.Name = "Form1";
            ((System.ComponentModel.ISupportInitialize)(this.pictureBoxPlateau)).EndInit();
            this.menuStripMain.ResumeLayout(false);
            this.menuStripMain.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Button button1;
        private System.Windows.Forms.PictureBox pictureBoxPlateau;
        private System.Windows.Forms.MenuStrip menuStripMain;
        private System.Windows.Forms.ToolStripMenuItem viewToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem textureToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem toolStripMenuItem2;
        private System.Windows.Forms.ToolStripMenuItem heightToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem plateauToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem unknow1ToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem unknow2ToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem passibleToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem unknow4ToolStripMenuItem;
        private System.Windows.Forms.TextBox textBoxGrid;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.ComboBox comboBoxLayerNumber;
        private System.Windows.Forms.ToolStripMenuItem fileToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem openToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem saveToolStripMenuItem;
        private System.Windows.Forms.ToolStripMenuItem newToolStripMenuItem;
        private System.Windows.Forms.Label label1;
    }
}

