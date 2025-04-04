import axios from 'axios';

export const VIEWS = {
  promoterDB: {
    title: 'Promoter DB',
    linkTo: 'NeuroSCAN',
    linkToRoute: 'https://neuroscan.net',
  },
  neuroScan: {
    title: 'NeuroSCAN',
    linkTo: 'Promoter DB',
    linkToRoute: 'https://promoters.wormguides.org/',
  },
};

export const VIEWERS = Object.freeze({
  InstanceViewer: 'Viewer',
  CphateViewer: 'Cphate',
});

export const ABOUT_CONTENT = [
  `is an initiative from the Yale University for Neurosciences, in partnership
with MetaCell and Bilte Co.`,
];

export const NEUROSCAN_ABOUT = [
  '<p>NeuroSCAN is a resource for exploring neuronal relationships and structures within the <em>C. elegans</em> nerve ring and across developmental stages. CPHATE graphs enable users to investigate how neuronal relationships evolve over time based on their contactomic profiles.</p>',
  '<p>To ensure meaningful comparative connectomics, we have collated segmented electron microscopy datasets across multiple developmental time points. This process required standardization and alignment to make the datasets truly comparable, allowing for accurate 3D visualization of neuron morphologies, neuron-neuron contact sites, and synaptic connections.</p>',
  '<p>By integrating these standardized datasets, NeuroSCAN provides a platform for exploring the developmental dynamics of the <em>C. elegans</em> nervous system.</p>',
  '<p>See <a href=\'https://elifesciences.org/reviewed-preprints/103977v1\' target=\'_blank\'>Koonce and Emerson et al.</a> for more details.</p>',
  '<p>To contribute to NeuroSCAN, contact <a href=\'mailto:daniel.colon-ramos@yale.edu\'>daniel.colon-ramos@yale.edu</a> / <a href=\'mailto:wmohler@neuron.uchc.edu\'>wmohler@neuron.uchc.edu</a></p>',
  '<h3>Data availability:</h3>',
  `<p>The data generated for NeuroSCAN is available in .OBJ file format and can be visualized locally using <strong>CytoSHOW</strong>. To explore the data, you can:<ul>
  <li>Use CytoSHOW: A program designed for interactive visualization.</li>
  <li>Access the data: Available at http://neuroscan.cytoshow.org/</li></ul></p>`,
  `<p>This ensures client-side exploration of neuronal structures and relationships.<ul>
  <li>For interactive contactome spreadsheets, see <a href="https://elifesciences.org/reviewed-preprints/103977v1" target="_blank">Koonce and Emerson et al.</a>, Supplementary Tables 8–13. These tables provide detailed insights into neuronal adjacencies and can be used alongside NeuroSCAN for comprehensive contactomic analysis.</li></ul></p>`,
  '<h3>CPHATE:</h3>',
  `<p><ul>
  <li>Brugnone <em>et al., Proc. IEEE Int. Conf. Big Data</em> (2019).<a href="https://arxiv.org/abs/1907.04463" target="_blank">DOI:10.1109/BigData47090.2019.9006013</a>.</li>
  <li>Moyle <em>et al., Nature</em> (2021). <a href="https://www.nature.com/articles/s41586-020-03169-5" target="_blank">DOI: 10.1038/s41586-020-03169-5</a></li>
  </ul></p>`,
  '<h3>Datasets used:</h3>',
  `<p><ul>
  <li><strong>L1, L2, L3, Adult (45h)</strong>: Witvliet <em>et al., Nature</em> (2021). <a href="https://www.nature.com/articles/s41586-021-03778-8" target="_blank">DOI: 10.1038/s41586-021-03778-8</a></li>
  <li><strong>L4, Adult (48h)</strong>: White <em>et al., Phil. Trans. R. Soc. Lond. B</em> (1986). <a href="https://royalsocietypublishing.org/doi/10.1098/rstb.1986.0056" target="_blank">DOI: 10.1098/rstb.1986.0056</a>; Cook <em>et al., Nature</em> (2019). <a href="https://www.nature.com/articles/s41586-019-1352-7" target="_blank">DOI: 10.1038/s41586-019-1352-7</a>.</li>
  </ul></p>`,
];

export const PROMOTERDB_ABOUT = [
  `The promoter database aims to share promoter expression data to visualize all C.
elegans neurons as they develop in the embryo. We use fluorescent membrane labels
driven by sparsely expressed promoters to see details of neural development, usually at
subcellular resolution. After imaging expression details on the diSPIM, we characterize
the expression with cell lineaging to identify cells that are labeled.`,
  `We invite the community to submit promoters for characterization using the “Suggest a
promoter” button in the top right. Promoters submitted to the WormGUIDES consortium
will be imaged using diSPIM, lineaged and the expression identity will be shared with
the community via our promoter database.`,
  'See resources about the data acquisition at:',
  `Duncan, L. H., Moyle, M. W., Shao, L., Sengupta, T., Ikegami, R., Kumar, A., Guo, M.,
Christensen, R., Santella, A., Bao, Z., Shroff, H., Mohler, W., Colón-Ramos, D. A.
Isotropic Light-Sheet Microscopy and Automated Cell Lineage Analyses to Catalogue
Caenorhabditis elegans Embryogenesis with Subcellular Resolution. &lt;em&gt;J. Vis.
Exp.&lt;/em&gt; (148), e59533, doi:10.3791/59533 (2019).`,
  `Kumar A, Wu Y, Christensen R, et al. Dual-view plane illumination microscopy for rapid
and spatially isotropic imaging. Nat Protoc. 2014;9(11):2555-2573.
doi:10.1038/nprot.2014.172`,
  `Also see our neurodevelopmental atlas for exploration of segmented neurons at:
wormguides.org`,
];

export const maxRecordsPerFetch = 30;

export const backendURL = process.env.REACT_APP_BACKEND_URL || '';
export const backendClient = axios.create({
  baseURL: backendURL,
});

export const VIEWER_MENU = {
  devStage: 'devStages',
  layers: 'layers',
  download: 'download',
  colorPicker: 'colorPicker',
};

export const filesURL = '/files';

export const NEURON_TYPE = 'neuron';
export const CONTACT_TYPE = 'contact';
export const SYNAPSE_TYPE = 'synapse';
export const NERVE_RING_TYPE = 'nervering';
export const CPHATE_TYPE = 'cluster';
export const SCALE_TYPE = 'scale';

export const CANVAS_BACKGROUND_COLOR_LIGHT = 0xffffff;
export const CANVAS_BACKGROUND_COLOR_DARK = 0x2c2c2c;

export const DOWNLOAD_SCREENSHOT = 'screenshot';
export const DOWNLOAD_OBJS = 'objects';

export const PROMOTER_MEDIA_TYPES = {
  video: 'video',
};

export const MAIL_SUGGEST_PROMOTER_TO = 'postmaster@wormguides.org';
export const MAIL_SUGGEST_PROMOTER_SUBJECT = 'Suggest a promoter';
export const MAIL_SUGGEST_PROMOTER_BODY = `Hi WormGUIDES team,\nOur group has found this promoter useful in our studies.
  We call the promoter (promoter name) and the strain name and/or primers for the promoter are as follows: The promoter
  has expression from (starting timepoint) to (ending timepoint). We see expression in these cells: \n We've attached an image of
  the promoter. \n\nThank you`;

export const MAIL_CONTACT_TO = 'support@wormguides.org';
export const MAIL_CONTACT_SUBJECT = '';
export const MAIL_CONTACT_BODY = '';

export const CANVAS_STARTED = 'STARTED';
export const CANVAS_FINISHED = 'FINISHED';

export const GREY_OUT_MESH_COLOR = {
  r: 128 / 255,
  g: 128 / 255,
  b: 128 / 255,
  a: 0.5,
};
